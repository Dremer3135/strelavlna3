<script lang="ts">
    import { currentState } from "$lib/stores/state";
    import type { Prob } from "$lib/types";
    import { isProbSolved } from "$lib/utils";

    let { prob, onSelect }: { prob: Prob, onSelect?: () => void } = $props();

    function handleClick() {
        if (onSelect) {
            onSelect();
        }
    }
</script>

<button class:selected={prob.focusedBy.includes($currentState.myId)} on:click={handleClick} class:solved={isProbSolved(prob)}>
    <h2 class="name">{prob.name}</h2>
    <div class="right">
        <h2 class="diff">[{prob.diff}]</h2>
        <div class="focus-indicator" class:focused={prob.focusedBy.length > (prob.focusedBy.includes($currentState.myId) ? 1 : 0)}></div>
    </div>
</button>

<style lang="scss">
    button {
        all: unset;
        cursor: pointer;
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: space-between;
        padding: 10px 20px;
        padding-right: 13px;
        border: 3px var(--color-purple) solid;
        border-radius: 5px;
        width: 100%;
        box-sizing: border-box;
        background-color: color-mix(in srgb, var(--color-purple) 5%, transparent 95%);

        &.selected {
            background-color: var(--color-purple);
            color: white;
        }

        h2 {
            margin: 0px;
        }

        &.solved {
            background-color: #F0F0F0BB;
            border-color: gray;
            transform: scale(0.96);
            
            .name {
                color: #555555;
            }

            .diff {
                color: #555555;
            }

            .focus-indicator {
                background-color: gray !important;
            }

            &.selected {
                background-color: gray;

                .diff {
                    color: white;
                }

                .focus-indicator {
                    background-color: white !important;
                }
            }
        }



        .name {
            font-family: 'Lexend';
            font-size: 22px;
            font-weight: 700;
            color: color-mix(in srgb, var(--color-purple) 30%, black 70%);
            // color: black
        }

        .right {
            display: flex;
            flex-direction: row;
            gap: 10px;
            align-items: center;
            justify-content: center;
            
            .diff {
                font-family: 'Lexend';
                font-size: 22px;
                font-weight: 500;
                transform: translateY(-1px);   
            }

            .focus-indicator {
                width: 15px;
                height: 15px;
                border-radius: 3px;
                background-color: var(--color-purple);

                opacity: 0;

                &.focused {
                    opacity: 1;
                }
            }
        }

        &.selected {
            .focus-indicator {
                background-color: white;
            }

            .name {
                color: white;
            }
        }

    }



</style>