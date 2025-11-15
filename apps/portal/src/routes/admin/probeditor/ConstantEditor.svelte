<script lang="ts">
    import type { ConstantsRecord, ConstantsResponse } from "$lib/pocketbase-types";
    import { editableConstants } from "$lib/stores/consts";
    import { createEventDispatcher } from "svelte";
    import ConstantGroupSelector from "./ConstantGroupSelector.svelte";
    import { getConstantEditedState } from "$lib/utils";

    let { values = $bindable(), groupNames, type, submit, discard }: { values: Partial<ConstantsRecord>, groupNames: string[], type: "add" | "edit" | "", submit: Function,  discard: Function } = $props();
    
    let sameVariableConstants: ConstantsResponse[] = $derived(Object.entries($editableConstants).filter(constant => getConstantEditedState(constant[1]).variable_name == values.variable_name && getConstantEditedState(constant[1]).id != values.id).map(constant => getConstantEditedState(constant[1])));
    let variableNameFree: boolean = $derived(sameVariableConstants.length == 0);
    let criteriaMet: boolean = $derived(
        variableNameFree &&
        values.group != ""
    );
    
    let notFreeTooltipActive: boolean = $state(false);


</script>

<!-- {#if values.group} -->
<main>
    <h2 class="title">{ type == "add" ? "Add a new constant" : "Edit constant" }</h2>
    <div class="content">
        <ConstantGroupSelector bind:value={ values.group } groups={ groupNames } />
        <input type="text" class="name" bind:value={values.name} placeholder="name">
        <input type="text" class="variable-name" bind:value={values.variable_name} placeholder="variable name">
        {#if !variableNameFree}
        <div class="variable-not-free-wrapper">
            <p>This name is already used!</p>
            <div class="anchor">
                <button onclick={() => { notFreeTooltipActive = !notFreeTooltipActive; }}>
                    <i class="fa-solid fa-circle-info"></i>
                </button>
                <div class="used-names-wrapper" class:active={ notFreeTooltipActive }>
                    {#each sameVariableConstants as constants}
                        <p>{constants.name}</p>
                    {/each}
                </div>
            </div>
        </div>
        {/if}
        <input type="text" class="description" bind:value={values.desc} placeholder="description">
        <input type="text" class="symbol" bind:value={values.symbol} placeholder="symbol">
        <input type="text" class="value" bind:value={values.value} placeholder="value">
        <input type="text" class="unit" bind:value={values.unit} placeholder="unit">
    </div>
    <div class="controls">
        <button class="discard" onclick={() => { discard(); }}>
            {#if type == "add"}
                <i class="fa-solid fa-xmark"></i>
            {:else}
                <i class="fa-solid fa-rotate-left"></i>
            {/if}
            {type == "add" ? "Dont add" : "Dont apply"}
        </button>
        <button class="submit" onclick={() => { if (criteriaMet) { submit(); }}} class:disabled={!criteriaMet}>
            <i class="fa-solid fa-check"></i>
            {type == "add" ? "Add constant" : "Apply changes"}
        </button>
    </div>
</main>
<!-- {/if} -->


<style lang="scss">
    main {
        display: flex;
        flex-direction: column;
        align-items: center;
        max-width: 600px;
        background-color: white;
        padding: 20px;
        border-radius: 5px;
        box-shadow: 0px 0px 20px 5px #00000055;
        gap: 20px;
        animation: spawn-animation cubic-bezier(0.215, 0.610, 0.355, 1) 0.2s forwards;

        .title {
            margin: 0px;
            font-family: 'Fredoka';
            font-weight: 600;
            font-size: 30px;
        }

        .content {
            display: flex;
            align-items: flex-start;
            justify-content: center;
            flex-direction: column;
            gap: 12px;

            input {
                all: unset;
                cursor: text;
                font-family: 'Fredoka';
                border: 3px lightgray solid;
                padding: 4px 12px;
                border-radius: 3px;
                background-color: #FAFAFA;
                width: 450px;

                transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;

                &.name:focus {
                    border-color: var(--color-orange);
                    background-color: color-mix(in srgb, var(--color-orange) 5%, transparent 95%);
                }
                &.variable-name:focus {
                    border-color: var(--color-orange);
                    background-color: color-mix(in srgb, var(--color-orange) 5%, transparent 95%);
                }
                &.description:focus {
                    border-color: var(--color-pink);
                    background-color: color-mix(in srgb, var(--color-pink) 5%, transparent 95%);
                }
                &.symbol:focus {
                    border-color: var(--color-purple);
                    background-color: color-mix(in srgb, var(--color-pink) 5%, transparent 95%);
                }
                &.value:focus {
                    border-color: var(--color-purple);
                    background-color: color-mix(in srgb, var(--color-purple) 5%, transparent 95%);
                }
                &.unit:focus {
                    border-color: var(--color-purple);
                    background-color: color-mix(in srgb, var(--color-purple) 5%, transparent 95%);
                }
            }

            .variable-not-free-wrapper {
                display: flex;
                align-items: center;
                justify-content: flex-start;
                gap: 0px;
                padding-bottom: 15px;

                p {
                    font-family: 'Fredoka';
                    font-size: 16px;
                    color: var(--color-pink);
                    margin: 0px;
                }

                .anchor {
                    position: relative;
                    
                    button {
                        all: unset;
                        cursor: pointer;
                        padding: 7px;
                        z-index: 10;

                        i {
                            font-size: 12px;
                            color: #666666;
                        }
                    }

                    .used-names-wrapper {
                        position: absolute;
                        top: 10px;
                        left: 50%;
                        transform: translateX(calc(-50%));
                        min-width: 300px;
                        display: flex;
                        flex-direction: column;
                        align-items: flex-start;
                        gap: 10px;
                        background-color: #FAFAFA;
                        border-radius: 5px;
                        box-shadow: 0px 0px 10px 5px #00000033;
                        padding: 20px;
                        box-sizing: border-box;
                        transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.2s;
                        opacity: 0;
                        z-index: -1;

                        &.active {
                            top: 40px;
                            left: 50%;
                            transform: translateX(calc(-50%));
                            opacity: 1;
                            z-index: 0;
                        }

                        p {
                            margin: none;
                            font-family: 'Fredoka';
                            font-size: 15px;
                        }
                    }
                }
            }
        }

        .controls {
            display: flex;
            justify-content: center;
            align-items: center;
            gap: 20px;

            button {
                all: unset;
                font-family: 'Fredoka';
                font-weight: 500;
                color: #333333;
                font-size: 16px;
                padding: 4px 12px;
                transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;

                &.submit {
                    border: 3px var(--color-yellow) solid;
                    background-color: color-mix(in srgb, var(--color-yellow) 2%, transparent 98%);
                    border-radius: 3px;
                    cursor: pointer;

                    &:hover {
                        background-color: color-mix(in srgb, var(--color-yellow) 10%, transparent 90%);

                    }
                }
            
                &.discard {
                    border: 3px var(--color-pink) solid;
                    background-color: color-mix(in srgb, var(--color-pink) 2%, transparent 98%);
                    border-radius: 3px;
                    cursor: pointer;

                    &:hover {
                        background-color: color-mix(in srgb, var(--color-pink) 10%, transparent 90%);

                    }
                }

                &.disabled {
                    cursor: not-allowed !important;
                    background-color: #FAFAFA !important;
                    border-color: lightgray !important;
                    color: #555555 !important;
                }
            
            }
        }
    }



    @keyframes spawn-animation {
        0% {
            transform: translateY(-50px);
            opacity: 0;
        }
        100% {
            transform: translateY(0px);
            opacity: 1;
        }
    }

    @keyframes despawn-animation {
        0% {
            transform: translateY(0px);
            opacity: 1;
        }
        100% {
            transform: translateY(-50px);
            opacity: 0;
        }
    }

    





</style>