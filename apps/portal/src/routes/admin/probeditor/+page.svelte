<script lang="ts">
    import ProbBanner from "$lib/components/admin/ProbBanner.svelte";
    import type { CorrectorsResponse, ProbsResponse, TypedPocketBase } from "$lib/pocketbase-types";
    // import { createPocketbaseInstance } from "$lib/server/pocketbase.js";
    import { onMount } from "svelte";
    import { probs } from "$lib/stores/probs.js";
    import { pb } from "$lib/pocketbase";
    let { data, form } = $props();


    let selectedProbId: string = $state("");
    let selectedProb: ProbsResponse | undefined = $state();
    let selectedProbName: string = $state("");
    let selectedProbText: string = $state("");
    let selectedProbAnswer: string = $state("");


    $effect(() => {
        selectedProb = $probs.find(prob => prob.id === selectedProbId);
    });

    $effect(() => {
        selectedProbName = selectedProb?.name ?? "";
        selectedProbText = selectedProb?.text ?? "";
        selectedProbAnswer = selectedProb?.answer ?? "";
    });

    $effect(() => {
        probs.update(currentProbs => currentProbs.map(prob => { 
            if (prob.id !== selectedProbId) { return prob; }
            else {
                prob.name = selectedProbName;
                return prob;
            }
        }));
    });

    $effect(() => {
        probs.update(currentProbs => currentProbs.map(prob => { 
            if (prob.id !== selectedProbId) { return prob; }
            else {
                prob.text = selectedProbText;
                return prob;
            }
        }));
    });

    $effect(() => {
        probs.update(currentProbs => currentProbs.map(prob => { 
            if (prob.id !== selectedProbId) { return prob; }
            else {
                prob.answer = selectedProbAnswer;
                return prob;
            }
        }));
    });


    let filters = $state([
        [
            {name: "all", function: (probs: ProbsResponse[]) => { return probs; }},
            {name: "my", function: (probs: ProbsResponse[]) => { return probs.filter(prob => prob.author == data.user?.id); }},
        ],[
            {name: "all", function: (probs: ProbsResponse[]) => { return probs; }},
            {name: "[A]", function: (probs: ProbsResponse[]) => { return probs.filter(prob => prob.diff == "A"); }},
            {name: "[B]", function: (probs: ProbsResponse[]) => { return probs.filter(prob => prob.diff == "B"); }},
            {name: "[C]", function: (probs: ProbsResponse[]) => { return probs.filter(prob => prob.diff == "C"); }},
        ]
    ]);

    let filteresSelected = $state(Array(filters.length).fill(0));
    let filteredProbs = $state($probs);


    $effect(() => {
        let temp = $probs;
        for (let i = 0; i < filters.length; i++) {
            temp = filters[i][filteresSelected[i]].function(temp);
        }
        filteredProbs = temp;
    });

    let filtersOpen = $state(false);

    async function addProb() {
        let prob: ProbsResponse;

        try {
            prob = await pb.collection("probs").create({
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
        }
        
        probs.update(currentProbs => [prob, ...currentProbs])
    }

    async function deleteProb(id: string) {
        await pb.collection("probs").delete(id);
        probs.update(currentProbs => currentProbs.filter(prob => prob.id != id));
    }

    async function saveChanges() {
        if (!selectedProb) return;

        try {
            await pb.collection('probs').update(selectedProb.id, {
                name: selectedProbName,
                text: selectedProbText,
                answer: selectedProbAnswer
            });
            // alert('Changes saved successfully!');
        } catch (err) {
            console.error('Failed to save changes:', err);
            alert('Failed to save changes.');
        }
    }
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
                {#each filteredProbs as prob}
                <button class="banner-select" onclick={() => {
                    saveChanges();
                    selectedProbId = prob.id;
                    }}>
                    <ProbBanner prob={prob} user={data.user as CorrectorsResponse} selected={ selectedProbId == prob.id } />
                </button>
                {/each}
            </div>
        </div>
        <div class="controls-wrapper">
            <button class="add" onclick={() => { addProb(); }}>Add</button>
            <button class="remove" onclick={() => { if (selectedProbId) { deleteProb(selectedProbId); }}}>Remove</button>
        </div>
    </div>
    <div class="main-content">
        <div class="left">
            <input type="text" class="title" placeholder="Title"
            value={selectedProbName}
            oninput={(event) => { selectedProbName = event.currentTarget.value; }}
            >
            <textarea name="text" id="prob-text" class="text"
            oninput={(event) => { selectedProbText = event.currentTarget.value; }}
            >{selectedProbText}</textarea>
            <input type="text" class="answer" placeholder="Answer"
            value={selectedProbAnswer}
            oninput={(event) => { selectedProbAnswer = event.currentTarget.value; }}
            >
        </div>
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