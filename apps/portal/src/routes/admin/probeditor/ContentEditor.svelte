<script lang="ts">
    import CodeEditor from "$lib/components/admin/CodeEditor.svelte";
    import LoadingAnimation from "$lib/components/general/LoadingAnimation.svelte";
    import { pocketbase } from "$lib/pocketbase";
    import type { ConstantsRecord, ProbsResponse } from "$lib/pocketbase-types";
    import { createEventDispatcher } from "svelte";
    import { editableConstants } from "$lib/stores/consts";
    import { getConstantEditedState } from "$lib/utils";
    import ConstantEditor from "./ConstantEditor.svelte";
    import { ProbContent, Latex, type ProbContentType } from 'shared'
    import SubmitButton from "$lib/components/general/SubmitButton.svelte";

    let { probRecord, value_name, value_text, value_answer, value_code, value_images, changesSaved }: { probRecord: ProbsResponse, value_name: string, value_text: string, value_answer: string, value_code: string, value_images: string[], changesSaved: boolean } = $props();

    const dispatch = createEventDispatcher<{ 
        "name": { value: string };
        "text": { value: string };
        "answer": { value: string };
        "code": { value: string };
        "image-add": { value: FileList };
        "image-delete": { value: string[] };
        "constant-add": { value: Partial<ConstantsRecord> },
        "constant-delete": { value: string[] }
    }>();

    let groupedConstantIds = $derived(Object.fromEntries(Object.entries(
        Object.groupBy(
            Object.entries($editableConstants),
            (item) => getConstantEditedState(item[1]).group ?? "Nezarazeno"
        )
    ).flatMap(([group, probs]) => probs === undefined ? [] : [[group, probs.map((e) => e[0])]])));

    let ConstantGroupNames = $derived(Object.keys(groupedConstantIds));

    let editingConstantValues: Partial<ConstantsRecord> = $state({});

    $effect(() => {
        if (selectedConstIds.length == 1 && $editableConstants[selectedConstIds[0]]) {
            editingConstantValues = getConstantEditedState($editableConstants[selectedConstIds[0]]);
        }
    });


    function handleInput(event: Event, field: "name" | "text" | "answer" | "code") {

        const target = event.target as HTMLInputElement | HTMLTextAreaElement;
        dispatch(field, { value: target.value });
    }

    let selectedTool = $state(0);


    let imageInput: HTMLInputElement;

    function handleImageUpload() {
        if (imageInput && imageInput.files && imageInput.files.length > 0) {
            dispatch("image-add", { value: imageInput.files });
        }
    }

    let revisionContent: ProbContentType = $state({});
    $effect(() => {
        probRecord;
        revisionContent = {};
    });
    let isRevisionLoading = $state(false);

    let selectedImages: number[] = $state([]);
    let imagesLoaded: Record<string, boolean> = $state({});
    let selectedConstIds: string[] = $state([]);
    let constantEditType: "add" | "edit" | "" = $state("");

    function constantEditHandleSubmit() {
        if (constantEditType == "add") {
            dispatch("constant-add", { value: editingConstantValues });

        } else {
            if (selectedConstIds.length == 1) {
                editableConstants.update(currentConstants => {
                    currentConstants[selectedConstIds[0]].edit = editingConstantValues;
                    return currentConstants;
                });
            }
        }

        constantEditType = "";
        selectedConstIds = [];
    }

    function constantEditHandleDiscard() {
        constantEditType = "";
        selectedConstIds = [];
    }


    function handleImageDelete() {
        if (selectedImages.length > 0) {
            let imageStrings: string[] = [];
            for(let idx of selectedImages) {
                imageStrings.push(value_images[idx]);
            }
            dispatch("image-delete", { value: imageStrings });
        }

        selectedImages = [];
    }

    $effect(() => {
        value_images;
        selectedImages = [];
    });

    $effect(() => {
        selectedTool;
        selectedConstIds = [];
    });

    async function reloadRevision() {
        isRevisionLoading = true;
        let response;
        try {
           response = await pocketbase.send("/api/code", { body: { id: probRecord.id }, method: 'post'});
            
            revisionContent = {
                name: response.name,
                text: response.text,
                answer: response.answer,
                images: response.images,
                diff: response.diff,
                id: response.id
            }

        } catch (err: any) {
            console.warn(err);
            console.log(response);
        } finally {
            isRevisionLoading = false;
        }
        
    }

</script>

<main>
    <div class="toolbar">
        <button class="item" class:selected={ selectedTool == 0 } onclick={() => { selectedTool = 0; }}>
            <i class="fa-regular fa-file-lines"></i>
            <p>text</p>
        </button>
        <button class="item" class:selected={ selectedTool == 1 } onclick={() => { selectedTool = 1; }}>
            <i class="fa-solid fa-code"></i>
            <p>code</p>
        </button>
        <button class="item" class:selected={ selectedTool == 2 } onclick={() => { selectedTool = 2; }}>
            <i class="fa-solid fa-image"></i>
            <p>images</p>
        </button>
        <button class="item" class:selected={ selectedTool == 3 } onclick={() => { selectedTool = 3; }}>
            <i class="fa-solid fa-table-list"></i>
            <p>constants</p>
        </button>
        <button class="item" class:selected={ selectedTool == 4 } onclick={() => { selectedTool = 4; }}>
            <i class="fa-solid fa-circle-check"></i>
            <p>revision</p>
        </button>
    </div>
    <div class="content">
        {#if selectedTool == 0}
        <div class="content-text">
            <input type="text" class="name" placeholder="Name"
                value={value_name}
                oninput={(e) => handleInput(e, "name")}
            >
            <textarea name="text" id="prob-text" class="text"
                oninput={(e) => handleInput(e, "text")}
            >{value_text}</textarea>
            <input type="text" class="answer" placeholder="Answer"
                value={value_answer}
                oninput={(e) => handleInput(e, "answer")}
            >
        </div>
        {:else if selectedTool == 1}
        <div class="content-code">
            <CodeEditor value={ value_code } on:code={(e) => { dispatch("code", e.detail) }}></CodeEditor>
        </div>
        {:else if selectedTool == 2}
        <div class="content-images">
            <div class="input">
                <input type="file" name="image-input" accept="image/*" bind:this={imageInput} onchange={handleImageUpload} class="hidden-input" id="image-input" multiple>
                <label for="image-input" class="image-input-button">
                    <i class="fa-regular fa-square-plus"></i>
                    <p>Upload images</p>
                </label>
                <button class="image-delete-button" class:inactive={ selectedImages.length == 0 } onclick={() => { handleImageDelete(); }}>
                    <i class="fa-regular fa-trash-can"></i>
                    <p>Delete { selectedImages.length == value_images.length && selectedImages.length > 1 ? "all" : selectedImages.length > 0 ? selectedImages.length : "" } image{ selectedImages.length > 1 ? 's' : '' }</p>
                </button>
            </div>

            <div class="images-holder">
                {#each value_images as image, i}
                    <button
                        class:selected={ selectedImages.includes(i) }
                        onclick={() => {
                            if (selectedImages.includes(i)) { selectedImages = selectedImages.filter((item) => item != i); }
                            else { selectedImages.push(i); }
                        }}
                    >
                        {#if !imagesLoaded[image]}
                            <LoadingAnimation />
                        {/if}
                        <img 
                            src={pocketbase.files.getURL(probRecord, image)}
                            alt={`Image ${i + 1}`}
                            class:hidden={!imagesLoaded[image]}
                            onload={() => imagesLoaded[image] = true}
                            onerror={() => imagesLoaded[image] = true}
                        >
                    </button>
                {/each}
            </div>
        </div>
        {:else if selectedTool == 3}
        <div class="content-constants">
            <div class="edit-dialog-wrapper" class:active={ constantEditType !== "" }>
                {#if constantEditType !== "" }
                <ConstantEditor bind:values={ editingConstantValues } groupNames={ ConstantGroupNames } type={ constantEditType }
                submit={ constantEditHandleSubmit }
                discard={ constantEditHandleDiscard }
                />
                {/if}
            </div>
            <div class="controls">
                <button class="add" onclick={() => {
                    constantEditType = "add"; 
                    editingConstantValues = {
                        name: "Název",
                        variable_name: "var_name",
                        desc: "Popis",
                        symbol: "Symbol",
                        value: 67,
                        unit: "Jednotka"
                    }
                    }}>
                    <i class="fa-regular fa-square-plus"></i>
                    <p>Add constant</p>
                </button>
                <button class="delete" class:inactive={ selectedConstIds.length == 0 } onclick={() => {
                    dispatch("constant-delete", { value: selectedConstIds });
                    selectedConstIds = [];
                    }}>
                    <i class="fa-regular fa-trash-can"></i>
                    <p>Delete {selectedConstIds.length == Object.entries($editableConstants).length ? "all " : selectedConstIds.length > 0 ? selectedConstIds.length + " " : "" }constant{selectedConstIds.length > 1 ? "s": "" }</p>
                </button>
                <button class="edit" class:inactive={ selectedConstIds.length != 1 } onclick={() => { if (selectedConstIds.length === 1) { constantEditType = "edit"; }}}>
                    <i class="fa-solid fa-pencil"></i>
                    <p>Edit constant</p>
                </button>
            </div>
            <div class="content">
                {#each Object.entries(groupedConstantIds) as idGroup, i }
                    <div class="table-wrapper" class:yellow={ i%2 == 0 } class:orange={ i%2 == 1 }>
                        <h2 class="title">{idGroup[0]}</h2>
                        <table>
                            <thead> 
                                <tr>
                                    <td>Name</td>
                                    <td>Description</td>
                                    <td>Symbol</td>
                                    <td>Value</td>
                                    <td>Unit</td>
                                </tr>
                            </thead>
                            <tbody>
                                {#each idGroup[1] as constantId}
                                    <tr onclick={() => {
                                        if (selectedConstIds.includes(constantId)) { selectedConstIds = selectedConstIds.filter(id => id != constantId); }
                                        else { selectedConstIds.push(constantId); }
                                    }}
                                    class:selected={ selectedConstIds.includes(constantId) }
                                    >
                                        <td class="name">
                                            <Latex text={getConstantEditedState($editableConstants[constantId]).name} />
                                        </td>
                                        <td class="descriptio">
                                            <Latex text={getConstantEditedState($editableConstants[constantId]).desc} />
                                        </td>
                                        <td class="symbol">
                                            <Latex text={getConstantEditedState($editableConstants[constantId]).symbol} />
                                        </td>
                                        <td class="value">
                                            <Latex text={getConstantEditedState($editableConstants[constantId]).value.toString()} />
                                        </td>
                                        <td class="unit">
                                            <Latex text={getConstantEditedState($editableConstants[constantId]).unit} />
                                        </td>
                                    </tr>
                                {/each}
                            </tbody>
                        </table>
                    </div>
                {/each} 
            </div>
        </div>
        {:else if selectedTool == 4}
        <div class="content-revision">
            <div class="controls">
                <button class="reload" onclick={() => { reloadRevision(); }} class:reloading={isRevisionLoading}>
                    Reload
                </button>
                {#if !changesSaved}
                <p>Unsaved changes will not apply!</p>
                {/if}
            </div>
            <div class="content">
                <ProbContent content={revisionContent}/>
            </div>
        </div>
        {/if}
    </div>
</main>

<style lang="scss">

    .toolbar {
        display: flex;
        flex-direction: row;
        justify-content: space-around;
        align-items: center;
    
        button.item {
            all: unset;
            cursor: pointer;
            padding: 10px 10px;
            display: flex;
            justify-content: center;
            align-items: center;
            gap: 8px;
            color: #555555;
            border-radius: 5px;
            background-color: transparent;
            box-sizing: border-box;

            &.selected {
                color: color-mix(in srgb, var(--color-pink) 40%, black 60%);
                border-bottom: 3px var(--color-pink) solid;
                background: linear-gradient(to bottom, transparent 20%, color-mix(in srgb, var(--color-pink) 10%, transparent 90%));
            }

            i {
                font-size: 20px;
            }

            p {
                font-family: 'Fredoka';
                font-size: 15px;
                margin: 0px;
            }
        }
    }


    main {
        display: flex;
        flex-direction: column;
        flex-grow: 1;
        // max-width: 750px;
        gap: 20px;
    }
    
    .content {
            flex-grow: 1;
            display: flex;
            flex-direction: column;
            min-height: 0px;
            overflow: hidden;
                    
        .content-text {
            display: flex;
            flex-direction: column;
            flex-grow: 1;
            gap: 10px;

            input, textarea {
                all: unset;
                padding: 5px 15px;
                border: 2px lightgray solid;
                border-radius: 3px;
                font-family: 'Lexend';
                font-size: 18px;
                font-weight: 700;
                background-color: color-mix(in srgb, lightgray 5%, transparent 95%);
                color: #333333;
                box-sizing: border-box;
                background-color: #F8F8F8;
                
                &:focus {
                    border-color: var(--color-pink);
                    background-color: color-mix(in srgb, var(--color-pink) 5%, transparent 95%);
                }
            }
        
            textarea {
                font-size: 16px;
                font-weight: 500;
                padding: 5px 8px;
                flex-grow: 1;
            }
        }

        .content-code {
            display: flex;
            flex-direction: column;
            flex-grow: 1;
            border-radius: 5px;
            overflow: hidden;
            border: 2px lightgray solid;
        }

        .content-images {
            display: flex;
            flex-direction: column;
            gap: 20px;
            min-height: 0px;

            .input {
                display: flex;
                flex-direction: row;
                align-items: center;
                justify-content: flex-start;
                gap: 15px;


                .image-input-button {
                    display: flex;
                    flex-direction: row;
                    align-items: center;
                    justify-content: center;
                    gap: 8px;
                    border: 3px var(--color-yellow) solid;
                    border-radius: 3px;
                    width: fit-content;
                    cursor: pointer;
                    padding: 5px 20px;
                    background-color: color-mix(in srgb, var(--color-yellow) 5%, transparent 95%);
                    color: #333333;
                    transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;

                    &:hover {
                        background-color: color-mix(in srgb, var(--color-yellow) 20%, transparent 80%);
                        color: color-mix(in srgb, var(--color-yellow) 20%, black 80%);
                    }
                    
                    i {
                        font-size: 20px;
                    }

                    p {
                        font-family: 'Fredoka';
                        font-size: 16px;
                        margin: 0px;
                    }
                }

                .image-delete-button {
                    display: flex;
                    flex-direction: row;
                    align-items: center;
                    justify-content: center;
                    gap: 8px;
                    border: 3px var(--color-pink) solid;
                    border-radius: 3px;
                    width: fit-content;
                    cursor: pointer;
                    padding: 5px 20px;
                    background-color: color-mix(in srgb, var(--color-pink) 5%, transparent 95%);
                    color: #333333;
                    transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;

                    &:hover {
                        background-color: color-mix(in srgb, var(--color-pink) 20%, transparent 80%);
                        color: color-mix(in srgb, var(--color-pink) 20%, black 80%);
                    }

                    &.inactive {
                        cursor: not-allowed;
                        color: #555555;
                        border-color: #AAAAAA;
                        background-color: color-mix(in srgb, #333333 5%, transparent 95%);
                    }
                    
                    i {
                        font-size: 20px;
                    }

                    p {
                        font-family: 'Fredoka';
                        font-size: 16px;
                        margin: 0px;
                    }

                    
                }

                .hidden-input {
                    display: none;
                }
            }

            .images-holder {
                display: flex;
                flex-direction: column;
                gap: 10px;
                // overflow: auto;

                .loading-placeholder {
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    min-height: 150px;
                    background-color: #f0f0f0;
                    border-radius: 5px;
                }

                .hidden {
                    display: none;
                }

                button {
                    all: unset;
                    max-width: 100%;
                    border-radius: 10px;
                    overflow: hidden;
                    border: 5px transparent solid;
                    transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;
                    cursor: pointer;
                    box-sizing: border-box;
                    width: fit-content;
                    padding: 5px;

                    
                    &.selected {
                        border-color: var(--color-pink);
                        border-style: solid;

                        img {
                            filter: sepia(100%) hue-rotate(280deg) saturate(200%);
                        }
                    }

                    &:hover {
                        // border-color: var(--color-pink);
                        transform: rotate(1.5deg);
                        box-shadow: 2px 10px 15px gray;

                    }

                    img {
                        max-width: 100%;
                        box-sizing: border-box;
                        border-radius: 5px;
                        overflow: hidden;
                        display: block;
                        transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.15s;
                    }

                }
            
            }
        }

        .content-constants {
            display: flex;
            flex-direction: column;
            gap: 20px;
            min-height: 0px;
            flex-grow: 1;

            .controls {
                display: flex;
                flex-direction: row;
                align-items: center;
                justify-content: flex-start;
                gap: 15px;

                button {
                    display: flex;
                    flex-direction: row;
                    align-items: center;
                    justify-content: center;
                    gap: 8px;
                    border-radius: 3px;
                    width: fit-content;
                    cursor: pointer;
                    padding: 5px 20px;
                    color: #333333;
                    transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;

                    i {
                        font-size: 20px;
                    }

                    p {
                        font-family: 'Fredoka';
                        font-size: 16px;
                        margin: 0px;
                    }
                }

                .add {
                    border: 3px var(--color-yellow) solid;
                    background-color: color-mix(in srgb, var(--color-yellow) 5%, transparent 95%);

                    &:hover {
                        background-color: color-mix(in srgb, var(--color-yellow) 20%, transparent 80%);
                        color: color-mix(in srgb, var(--color-yellow) 20%, black 80%);
                    }
                }

                .delete {
                    border: 3px var(--color-pink) solid;
                    background-color: color-mix(in srgb, var(--color-pink) 5%, transparent 95%);

                    &:hover {
                        background-color: color-mix(in srgb, var(--color-pink) 20%, transparent 80%);
                        color: color-mix(in srgb, var(--color-pink) 20%, black 80%);
                    }

                    &.inactive {
                        cursor: not-allowed;
                        color: #555555;
                        border-color: #AAAAAA;
                        background-color: color-mix(in srgb, #333333 5%, transparent 95%);
                    }
                    
                }

                .edit {
                    border: 3px var(--color-purple) solid;
                    background-color: color-mix(in srgb, var(--color-purple) 5%, transparent 95%);

                    &:hover {
                        background-color: color-mix(in srgb, var(--color-purple) 20%, transparent 80%);
                        color: color-mix(in srgb, var(--color-purple) 20%, black 80%);
                    }

                    &.inactive {
                        cursor: not-allowed;
                        color: #555555;
                        border-color: #AAAAAA;
                        background-color: color-mix(in srgb, #333333 5%, transparent 95%);
                    }
                    
                }
            }

            .content {
                overflow-y: auto;
                scrollbar-width: 0px;
                min-height: 0;
                display: flex;
                flex-direction: column;
                align-items: center;
                gap: 30px;
                flex-grow: 1;

                &::-webkit-scrollbar {
                    display: none;
                }
            
                .table-wrapper {
                    display: flex;
                    flex-direction: column;
                    align-items: center;
                    width: 100%;
                    background-color: #FCFCFC;

                    border: 3px lightgray solid;
                    border-radius: 5px;

                    padding: 10px;
                    padding-top: 0px;
                    
                    font-family: 'Lexend';

                    &.yellow {
                        background-color: color-mix(in srgb, var(--color-yellow) 3%, transparent 97%);
                        border-color: color-mix(in srgb, var(--color-yellow) 50%, white 50%);

                        tbody {
                            td {
                                background-color: #f0f0f0C0;
    
                            }
                        }
                    }

                    &.orange {
                        background-color: color-mix(in srgb, var(--color-orange) 3%, transparent 97%);
                        border-color: color-mix(in srgb, var(--color-orange) 70%, white 30%);

                        tbody {
                            td {
                                background-color: #f0f0f0C0;
    
                            }
                        }
                    }

                    .title {
                        margin: 20px;
                    }

                    table {
                        width: 100%;
                        
                        thead {
                            td {
                                text-align: center;
                                font-family: 'Fredoka';
                                font-weight: 600;
                                font-size: 18px;
                            }
                        }
    
                        tbody {
                            tr {
                                
                                td {
                                    padding: 3px 10px;
                                    border-radius: 3px;
                                    cursor: pointer;
                                    transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.15s;
                                    
                                    &.name {
                                        border-left: 3px solid transparent;
                                    }
                                }
        
                                &:hover {
                                    td {
                                        background-color: color-mix(in srgb, var(--color-pink) 20%, transparent 80%);
                                    }
                                }
    
                                &.selected {
                                    td {
                                        background-color: color-mix(in srgb, var(--color-pink) 10%, transparent 90%);
                                        
                                        &.name {
                                            border-left: 3px solid var(--color-pink);    
                                        }
                                    }
                                }
                            }
                        }

                    }


                }
            }

            .edit-dialog-wrapper {
                display: none;
                position: fixed;
                top: 0px;
                left: 0px;
                width: 100%;
                height: 100%;
                
                background-color: #00000000;
                backdrop-filter: blur(0.0px);
                z-index: -100;
                transition: background-color cubic-bezier(0.215, 0.610, 0.355, 1) 0.5s;
                    
                &.active {
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    background-color: #00000033;
                    backdrop-filter: blur(2.0px);
                    z-index: 100;
                    
                }
            }
        }

        .content-revision {
            display: flex;
            flex-direction: column;
            // gap: 20px;
            min-height: 0px;
            flex-grow: 1;
            
            .controls {
                width: 100px;
                align-self: flex-end;
                // padding-right: 50px;
                box-sizing: border-box;
                display: flex;
                flex-direction: column;
                gap: 10px;
                position: relative;

                .reload {
                    all: unset;
                    background-color: #444444;
                    cursor: pointer;
                    border-radius: 4px;
                    padding: 7px 15px;
                    transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;
                    font-family: 'Fredoka';
                    font-size: 15px;
                    font-weight: 600;
                    text-align: center;
                    color: white;
                    width: 100%;
                    box-sizing: border-box;
                    
                    &:hover {
                        background-color: black;
                    }

                    &.reloading {
                        background-color: #F0F0F0;
                    }
                }

                p {
                    position: absolute;
                    top: 40px;
                    margin: 0px;
                    font-size: 14px;
                    font-family: 'Fredoka';
                    color: var(--color-orange);
                    // width: fit-content;
                    // text-wrap: nowrap;
                }

                // box-sizing: border-box;
            }

        }
    }
</style>