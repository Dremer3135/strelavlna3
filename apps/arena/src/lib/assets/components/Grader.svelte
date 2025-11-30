<script lang="ts">
    import type { Prob } from "$lib/types";
    import { isProbSolved, isProbUngraded } from "$lib/utils";
    import Button from "./Button.svelte";

    let { prob, answer, onAccept, onReject }: { answer: string, prob: Prob | undefined, onAccept: () => void, onReject: () => void } = $props();
</script>

<svelte:body onkeydown={(e) => {
    if (!prob || isProbSolved(prob)) return;
    if (e.shiftKey) {
        if (e.key == "A") {
            onAccept();
        } else if (e.key == "R" && isProbUngraded(prob)){
            onReject();
        }
}}}/>

<main>
    <div class="answer">Odpověď týmu: <span class="bold">{answer}</span></div>
    <div class="separator"></div>
    <div class="actual">Správná odpověď: <span class="bold">{prob?.answer}</span></div>
    <div class="controls">
        <Button theme="yellow" disabled={!prob || isProbSolved(prob)} onclick={() => {
            if (!prob || isProbSolved(prob)) return;
            onAccept();
        }}>
            <p>Accept <span class="thin">(Shift + A)</span></p>
        </Button>
        <Button theme="pink" disabled={!prob || !isProbUngraded(prob) || isProbSolved(prob)} onclick={() => {
            if (!prob || !isProbUngraded(prob) || isProbSolved(prob)) return;
            onReject();
        }}>
            <p>Reject <span class="thin">(Shift + R)</span></p>
        </Button>
    </div>
</main>

<style lang="scss">
    main {
        width: 100%;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;

        .answer {
            width: 100%;
            padding: 10px;
            box-sizing: border-box;
            font-family: 'Lexend';
            font-size: 18px;
            font-weight: 500;
            color: color-mix(in srgb, var(--color-purple) 30%, black 70%);
            background-image: linear-gradient(to top, 
                color-mix(in srgb, var(--color-purple) 20%, transparent 80%) 0%,
                transparent 100%,
            );
        }

        .separator {
            width: 100%;
            border-top: 3px var(--color-purple) dashed;
        }

        .actual {
            width: 100%;
            padding: 10px;
            box-sizing: border-box;
            font-family: 'Lexend';
            font-size: 18px;
            font-weight: 500;
            background-image: linear-gradient(to bottom, 
                color-mix(in srgb, var(--color-yellow) 20%, transparent 80%) 0%,
                transparent 100%,
            );
        }

        .bold {
            font-weight: 700;
        }

        .controls {
            width: 100%;
            padding: 40px 0px;
            box-sizing: border-box;
            display: flex;
            justify-content: space-evenly;
            align-items: center;

            p {
                font-family: 'Lexend';
                font-size: 16px;
                font-weight: 600;
                margin: 0px;


                .thin {
                    font-size: 14px;
                    font-weight: 400;
                }
            }
        }

    }

</style>