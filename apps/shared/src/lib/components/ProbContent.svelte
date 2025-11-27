<script lang="ts">
    import Latex from "./Latex.svelte";
    import type { ProbContentType } from "$lib/types.ts";

    let { content, onCopy }: { content: ProbContentType, onCopy?: (text: string) => void } = $props()


    $effect(() => {
        console.log(content.answer);
    })
</script>

<main oncopy={(e) => {
        const selection = window.getSelection()?.toString();
        if (selection && onCopy) {
            onCopy(selection);
        }
    }}>
    <h1 class="name">{content.name}</h1>
    <p class="text"><Latex text={content.text}/></p>
    {#each content.images as image}
        <img class="image" src={image} alt="problem">
    {/each}
    {#if content.answer}
    <p class="answer"><Latex text={content.answer}/></p>
    {/if}
</main>

<style lang="scss">
    main {
        flex-grow: 1;
        display: flex;
        flex-direction: column;
        gap: 20px;
        justify-content: center;
        align-items: flex-start;
        width: 100%;
        box-sizing: border-box;

        .name {
            font-family: 'Lexend';
            font-size: 35px;
            font-weight: 600;
            margin: 0px;
        }

        .text {
            flex-grow: 1;
            text-wrap: wrap;
            font-family: 'Lexend';
            font-size: 16px;
            font-weight: 500;
        }

        .image {
            max-height: 400px;
            max-width: 600px;
            object-fit: contain;
        }

        .answer {
            font-family: 'Lexend';
            font-size: 20px;
            font-weight: 600;
        }
    }
</style>