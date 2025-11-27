<script lang="ts">
    import Button from "$lib/assets/components/Button.svelte";
    import ProbSelect from "$lib/assets/components/ProbSelect.svelte";
    import { probs } from "$lib/stores/probs";
    import { onMount } from "svelte";
    import type { MessageType } from "$lib/types";
    import { focusedProb, currentState } from "$lib/stores/state";
    import { ProbContent } from "shared"
    import AnswerInput from "$lib/assets/components/AnswerInput.svelte";
    import Chat from "$lib/assets/components/Chat.svelte";

    let { buy, sell, focus, chat }: { buy: (diff: string) => void, sell: (id: string) => void, focus: (id: string) => void, chat: (probId: string, message: MessageType) => void } = $props();

    const difficulties = ["A", "B", "C"];

    let buySelectOpened: boolean = $state(false);

    $effect(() => {
        console.log($focusedProb);
    })

</script>


<main>
    <div class="left-panel">
        <div class="team-stats">
            <h2 class="name">{$currentState.teamName}</h2>
            <h3 class="money">{$currentState.money} DC</h3>
        </div>
        <div class="prob-selector">
            {#each Object.values($probs) as prob}
                <ProbSelect prob={prob} onSelect={() => {
                    focus(prob.id);
                }}/>
            {/each}
        </div>
        
        <div class="controls">
            <div class="buy">
                <Button disabled={false} theme="yellow" onclick={() => {buySelectOpened = !buySelectOpened}}>
                    <i class="fa-solid fa-angle-down"></i>
                    <p class="buy">Koupit</p>
                </Button>
                <div class="dropdown" class:opened={buySelectOpened}>
                    {#each difficulties as diff, i}
                        <Button disabled={false} theme={ ["yellow", "orange", "pink", "purple"][i%4] as "yellow" | "orange" | "pink" | "purple" }
                        onclick={() => {
                            buySelectOpened = !buySelectOpened;
                            buy(diff);
                        }}
                        >
                            <p class="buy">[{diff}]</p>
                        </Button>
                        
                    {/each}
                </div>
            </div>
            <div class="sell">
                <Button disabled={false} theme="pink"
                onclick={() => {
                    if ($focusedProb) {
                        sell($focusedProb.id);
                    }
                }}
                >
                    <i class="fa-solid fa-trash-can"></i>
                    <p class="buy">Prodat</p>
                </Button>
            </div>
        </div>
    </div>

    <div class="content">
        {#if $focusedProb}
            <div class="prob">
                <ProbContent content={$focusedProb} onCopy={(text: string) => {
                    console.log("copying");
                    
                    if (!$focusedProb) return;

                    chat($focusedProb.id, {
                        type: "copy",
                        value: text,
                        origin: "sent"
                    })
                }}/>
                <AnswerInput submitAnswer={(answer: string) => {
                    if (!$focusedProb) return;

                    chat($focusedProb.id, {
                        type: "answer",
                        value: answer,
                        origin: "sent"
                    });
                }}
                onPaste={(text: string) => {
                    if (!$focusedProb) return;

                    chat($focusedProb.id, {
                        type: "paste",
                        value: text,
                        origin: "sent"
                    });
                }}
                />   

            </div>
            <div class="chat">
                <Chat prob={$focusedProb} send={(message) => {
                    if (!$focusedProb) return;

                    chat($focusedProb.id, message);
                }} /> 
            </div>
        {/if}
    </div>  

</main>

<style lang="scss">

    main {
        flex-grow: 1;
        display: flex;
        flex-direction: row;

        .left-panel {
            width: 350px;
            display: flex;
            flex-direction: column;

            .team-stats {
                display: flex;
                flex-direction: row;
                align-items: center;
                justify-content: space-between;
                padding: 20px;
                box-sizing: border-box;
                background-color: color-mix(in srgb, var(--color-purple) 5%, transparent 95%);
                border-right: 1px color-mix(in srgb, var(--color-purple) 30%, transparent 70%) solid;
                border-bottom: 1px color-mix(in srgb, var(--color-purple) 30%, transparent 70%) solid;
                gap: 10px;

                .name {
                    font-family: 'Lexend';
                    font-weight: 700;
                    font-size: 21px;
                    margin: 0px;
                    color: color-mix(in srgb, var(--color-purple) 30%, black 70%)
                }

                .money {
                    align-self: flex-start;
                    font-family: 'Lexend';
                    font-weight: 600;
                    font-size: 23px;
                    margin: 0px;
                    text-wrap: nowrap;
                    background-color: color-mix(in srgb, var(--color-purple) 5%, transparent 95%);
                    border: 2px color-mix(in srgb, var(--color-purple) 20%, transparent 80%) solid;
                    padding: 2px 7px;
                    border-radius: 3px;
                }
            }

            .prob-selector {
                display: flex;
                flex-direction: column;
                gap: 10px;
                overflow-y: auto;
                min-height: 0px;
                background-color: color-mix(in srgb, var(--color-purple) 5%, transparent 95%);
                padding: 20px;
                box-sizing: border-box;
                width: 100%;
                flex-grow: 1;
                border-right: 1px color-mix(in srgb, var(--color-purple) 30%, transparent 70%) solid;
                // border-right: 1px var(--color-purple) solid;
                
            }
            
            .controls {
                background-color: color-mix(in srgb, var(--color-purple) 5%, transparent 95%);
                border-top: 1px color-mix(in srgb, var(--color-purple) 30%, transparent 70%) solid;
                border-right: 1px color-mix(in srgb, var(--color-purple) 30%, transparent 70%) solid;
                padding: 20px 20px;
                display: flex;
                justify-content: space-between;
                align-items: start;

                p {
                    font-family: 'Lexend';
                    font-size: 16px;
                    font-weight: 600;
                    margin: 0px;
                }

                i {
                    color: black;
                }

                .dropdown {
                    opacity: 0;
                    padding: 10px;
                    box-sizing: border-box;
                    display: flex;
                    flex-direction: column;
                    justify-content: center;
                    align-items: flex-start;
                    gap: 10px;

                    &.opened {
                        opacity: 1;
                    }
                }

            }
        }
        
        .content {
            display: flex;
            flex-grow: 1;
            .prob {
                flex-grow: 1;
                display: flex;
                flex-direction: column;
                padding: 20px;
                box-sizing: border-box;

            }

            .chat {
                display: flex;
                flex-direction: column;
            }
        }
    }


</style>