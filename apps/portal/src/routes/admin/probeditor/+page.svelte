<script lang="ts">
    import ContentEditor from "./ContentEditor.svelte";
    import ProbBanner from "$lib/components/admin/ProbBanner.svelte";
    import type { ConstantsRecord, CorrectorsResponse, ProbsResponse, TypedPocketBase } from "$lib/types/pocketbase-types";
    // import { createPocketbaseInstance } from "$lib/server/pocketbase.js";
    import { onMount } from "svelte";
    import { editableProbs } from "$lib/stores/probs.js";
    import { editableConstants } from "$lib/stores/consts";
    import type { EditableConstant } from "$lib/types.js";
    import { pocketbase } from "$lib/pocketbase";
    import type { EditableProb } from "$lib/types.js";
    import { filterRecord, getProbEditedState, isProbEdited, isConstantEdited } from "$lib/utils.js";
    import { getRequestEvent } from "$app/server";
    let { data, form } = $props();



    let selectedProbId = $state<string | undefined>(undefined);
    let selectedProb = $derived(selectedProbId ? $editableProbs[selectedProbId] : undefined);



    let filters = $state([
        [
            {name: "all", function: (probIds: string[]) => { return probIds; }},
            {name: "my", function: (probIds: string[]) => { return probIds.filter( probId => getProbEditedState($editableProbs[probId]).author == data.user?.id )}},
            {name: "free", function: (probIds: string[]) => { return probIds.filter( probId => getProbEditedState($editableProbs[probId]).author == "" )}}
        ],[
            {name: "all", function: (probIds: string[]) => { return probIds; }},
            {name: "[A]", function: (probIds: string[]) => { return probIds.filter(probId => getProbEditedState($editableProbs[probId]).diff == "A"); }},
            {name: "[B]", function: (probIds: string[]) => { return probIds.filter(probId => getProbEditedState($editableProbs[probId]).diff == "B"); }},
            {name: "[C]", function: (probIds: string[]) => { return probIds.filter(probId => getProbEditedState($editableProbs[probId]).diff == "C"); }},
        ]
    ]);

    let filteresSelected = $state(Array(filters.length).fill(0));
    let filteredProbIds = $state(Object.keys($editableProbs));

    $effect(() => {
        let temp = Object.keys($editableProbs);
        for (let i = 0; i < filters.length; i++) {
            temp = filters[i][filteresSelected[i]].function(temp);
        }
        filteredProbIds = temp;
    });

    let filtersOpen = $state(false);

    async function addProb() {
        let probResponse: ProbsResponse;

        try {
            probResponse = await pocketbase.collection("probs").create({
                name: "Moje uloha",
                diff: "A",
                auto: false,    
                infinite: false,
                code: "",
                text: "Tohle bude text ulohy",
                answer: "Tohle je odpoved",
                contests: [],
                author: data.user?.id,
                queue: []
            });
        } catch (err) {
            console.log(err);
            return;
        }

        const newProb: EditableProb = {
            prob: probResponse,
            edit: {}
        }
        
        editableProbs.update(currentProbs => {
            currentProbs[newProb.prob.id] = newProb;
            return currentProbs;
        });

        selectedProbId = newProb.prob.id;
    }

    async function deleteProb(id: string) {
        await pocketbase.collection("probs").delete(id);
        editableProbs.update(currentProbs => {
            delete currentProbs[id];
            return currentProbs;
        });

        selectedProbId = undefined;
    }

    async function saveChanges() {
        const updatePromises = [];
        for (let probId of Object.keys($editableProbs)) {
            if (isProbEdited($editableProbs[probId])) {
                const promise = pocketbase.collection('probs').update(probId, {
                    ...$editableProbs[probId].edit
                });
                updatePromises.push(promise);
            }
        }

        for (let constId of Object.keys($editableConstants)) {
            if (isConstantEdited($editableConstants[constId])) {
                const promise = pocketbase.collection('constants').update(constId, {
                    ...$editableConstants[constId].edit
                });
                updatePromises.push(promise);
            }
        }

        try {
            await Promise.all(updatePromises);

            editableProbs.update(currentProbs => {
                for (let probId of Object.keys(currentProbs)) {
                    if (Object.keys(currentProbs[probId].edit).length > 0) {
                        Object.assign(currentProbs[probId].prob, currentProbs[probId].edit);
                        currentProbs[probId].edit = {};
                    }
                }
                return currentProbs;
            });

            editableConstants.update(currentConstants => {
                for (let constId of Object.keys(currentConstants)) {
                    if (Object.keys(currentConstants[constId].edit).length > 0) {
                        Object.assign(currentConstants[constId].constant, currentConstants[constId].edit);
                        currentConstants[constId].edit = {};
                    }
                }
                return currentConstants;
            });

        } catch (err) {
            console.error('Failed to save changes:', err);
            alert('Failed to save changes.');
        }
    }

    async function uploadImages(files: FileList) {
        if (!selectedProb) { return; }
        const imagesArray = Array.from(files);
        const updatedProb = await pocketbase.collection("probs").update(selectedProb.prob.id, {
            "images+": imagesArray
        });

        editableProbs.update(currentProbs => {
            currentProbs[updatedProb.id].prob.images = updatedProb.images;
            return currentProbs;
        });

        if (selectedProb && selectedProb.prob.id == updatedProb.id) {
            selectedProb = $editableProbs[selectedProb.prob.id];
        }
    }

    async function deleteImages(names: string[]) {
        if (!selectedProb) { return; }

        const updatedProb = await pocketbase.collection("probs").update(selectedProb.prob.id, {
            "images-": names
        });

        editableProbs.update(currentProbs => {
            currentProbs[updatedProb.id].prob.images = updatedProb.images;
            return currentProbs;
        });

        if (selectedProb && selectedProb.prob.id == updatedProb.id) {
            selectedProb = $editableProbs[selectedProb.prob.id];
        }
    }
    
    async function addConstant(values: Partial<ConstantsRecord>) {
        let newConstant = await pocketbase.collection("constants").create(values);

        editableConstants.update(currentConstants => {
            currentConstants[newConstant.id] = {
                constant: newConstant,
                edit: {}
            }
            return currentConstants;
        });
    }

    async function deleteConstants(ids: string[]) {
        const promisses = [];
        for (let id of ids) {
            promisses.push(pocketbase.collection("constants").delete(id));
        }

        await Promise.all(promisses);

        editableConstants.update(currentConstants => {
            currentConstants = Object.fromEntries(Object.entries(currentConstants).filter(constant => !ids.includes(constant[1].constant.id))) 
            return currentConstants;
        });
    }
    
    // $effect(() => {
    //     $editableProbs;

    //     console.log("updated");
    // });



</script>



<main>
    <div class="banners-wrapper">
        <div class="filters-wrapper">
            <button class="dropdow-button" onclick={() => {filtersOpen = !filtersOpen}}>
                <div class="left">
                    <i class="fa-solid fa-sliders"></i>
                    <p>Filters</p>
                </div>
                <div class="right">
                    <i class="fa-solid fa-angle-down" class:rotated={filtersOpen}></i>
                </div>
            </button>
            <div class="content" class:open={filtersOpen}>
                {#each filters as filter_row, i}
                    <div class="filter-row">
                        {#each filter_row as filter, j}
                            <button class="filter-button" class:selected={filteresSelected[i] == j}
                                onclick={() => {filteresSelected[i] = j}}
                            >
                                {filter.name}
                            </button>
                        {/each}
                    </div>
                {/each}
            </div>
        </div>
        <div class="banners-scrollview">
            <div class="banners-holder">
                {#each filteredProbIds as probId}
                <button class="banner-select" onclick={() => {
                    saveChanges();
                    selectedProbId = probId;
                    }}>
                    <ProbBanner eprob={ $editableProbs[probId] } user={data.user as CorrectorsResponse} selected={ selectedProb?.prob.id == $editableProbs[probId].prob.id } />
                </button>
                {/each}
            </div>
        </div>
        <div class="controls-wrapper">
            <button class="add" onclick={() => { addProb(); }}>Add</button>
            <button class="remove" onclick={() => { if (selectedProb) { deleteProb(selectedProb.prob.id); }}}>Remove</button>
        </div>
    </div>
    <div class="main-content">
        {#if selectedProb}
            <ContentEditor
                probRecord={getProbEditedState(selectedProb)}
                value_name={getProbEditedState(selectedProb).name}
                value_text={getProbEditedState(selectedProb).text}
                value_answer={getProbEditedState(selectedProb).answer}
                value_code={getProbEditedState(selectedProb).code}
                value_images={getProbEditedState(selectedProb).images}
                on:name={(e) => {
                    editableProbs.update(currentProbs => {
                        if (selectedProb) {
                            currentProbs[selectedProb.prob.id].edit.name = e.detail.value;
                        }
                        return currentProbs;
                    });
                }}
                on:text={(e) => {
                    editableProbs.update(currentProbs => {
                        if (selectedProb) {
                            currentProbs[selectedProb.prob.id].edit.text = e.detail.value;
                        }
                        return currentProbs;
                    });
                }}
                on:answer={(e) => {
                    editableProbs.update(currentProbs => {
                        if (selectedProb) {
                            currentProbs[selectedProb.prob.id].edit.answer = e.detail.value;
                        }
                        return currentProbs;
                    });
                }}
                on:code={(e) => {
                    editableProbs.update(currentProbs => {
                        if (selectedProb) {
                            currentProbs[selectedProb.prob.id].edit.code = e.detail.value;
                        }
                        return currentProbs;
                    });
                }}
                on:image-add={(e) => {
                    uploadImages(e.detail.value);
                }}
                on:image-delete={(e) => {
                    deleteImages(e.detail.value);
                }}
                on:constant-add={(e) => {
                    addConstant(e.detail.value);
                }}
                on:constant-delete={(e) => {
                    deleteConstants(e.detail.value);
                }}
            />
        {/if}
    </div>

</main>

<style lang="scss">
    main {
        flex: 1;
        padding: 0px;
        display: flex; /* Make main a flex container as well */
        flex-direction: row;
        min-height: 0;
    }

    .main-content {
        display: flex;
        flex-direction: row;
        flex-grow: 1;
        padding: 20px;
    }
    

    
    .filters-wrapper {
        border-bottom: 1px color-mix(in srgb, var(--color-purple) 20%, transparent 80%) solid;
        
        .content {
            display: flex;
            justify-content: flex-start;
            align-items: flex-start;
            flex-direction: column;
            gap: 20px;
            padding: 20px 20px;
            box-sizing: border-box;
            min-height: 0px;
            height: 0px;
            overflow: hidden;
            padding-top: 0px;
            padding-bottom: 0px;
            
            &.open {
                height: fit-content;
                padding-bottom: 20px;
            }
        
            .filter-row {
                display: flex;
                justify-content: flex-start;
                align-items: center;
                gap: 10px;
            
                .filter-button {
                    all: unset;
                    cursor: pointer;
                    font-family: 'Fredoka';
                    font-weight: 700;
                    border: 2px lightgray solid;
                    border-radius: 3px;
                    padding: 2px 15px;
                    background-color: transparent;                 
                    transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;   
                    
                    &.selected {
                        border-color: var(--color-pink);
                        background-color: color-mix(in srgb, var(--color-pink) 20%, transparent 80%);
                    }
                }
                
            }
        }


        .dropdow-button {
            all: unset;
            cursor: pointer;
            display: flex;
            justify-content: space-between;
            align-items: center;
            width: 100%;
            padding: 5px 20px;
            box-sizing: border-box;
            
            .left {
                display: flex;
                justify-content: flex-start;
                align-items: center;
                gap: 10px;
                
                i {
                    font-size: 20px;
                }
            }

            .right i {
                transform: rotate(0deg);
                transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;

                &.rotated {
                    transform: rotate(180deg);
                }
            }

            p {
                font-family: 'Fredoka';
                font-size: 16px;
                font-weight: 700;
            }
        }
    }

    .banners-wrapper {
        display: flex;
        flex-direction: column;
        width: fit-content;
        padding-bottom: 50px;
        padding-top: 15px;
        background-color: color-mix(in srgb, var(--color-purple) 3%, transparent 97%);
        border-right: color-mix(in srgb, var(--color-purple) 20%, transparent 80%) 1px solid;
        min-height: 0;
        height: 100%;
    }


    .banners-scrollview {
        overflow: auto;
        width: 100%;
        flex: 1;
        padding: 20px;
        box-sizing: border-box;
        min-height: 0;
        
    }
    .banners-scrollview::-webkit-scrollbar {
        display: none;
    }
    
    .banners-holder {
        display: flex;
        justify-content: flex-start;
        align-items: flex-start;
        flex-direction: column;
        box-sizing: border-box;
        gap: 10px;
    }

    button.banner-select {
        all: unset;
        cursor: pointer;
        box-sizing: border-box;
        width: 100%;
    }

    .controls-wrapper {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 15px;
        border-top: color-mix(in srgb, var(--color-purple) 20%, transparent 80%) 1px solid;
        padding-top: 30px;
    }

    button.add {
        all: unset;
        cursor: pointer;
        border: 3px var(--color-yellow) solid;
        background-color: color-mix(in srgb, var(--color-yellow) 5%, transparent 95%);
        padding: 5px 20px;
        border-radius: 3px;
        font-family: 'Lexend';
        font-weight: 500;
        transition: all 0.5s cubic-bezier(0.215, 0.610, 0.355, 1);
    }
    button.add:hover {
        background-color: color-mix(in srgb, var(--color-yellow) 20%, transparent 80%);
    }


    button.remove {
        all: unset;
        cursor: pointer;
        border: 3px var(--color-pink) solid;
        background-color: color-mix(in srgb, var(--color-pink) 5%, transparent 95%);
        text-align: center;
        padding: 5px 20px;
        border-radius: 3px;
        font-family: 'Lexend';
        font-weight: 500;
        transition: all 0.5s cubic-bezier(0.215, 0.610, 0.355, 1);
    }
    button.remove:hover {
        background-color: color-mix(in srgb, var(--color-pink) 20%, transparent 80%);
    }


</style>
