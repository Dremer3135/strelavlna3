<script lang="ts">
    let { text }: { text: string } = $props();

    let replaced: boolean[] = $state(Array(text.length).fill(false));



</script>

<main>
    {#each text as char, i}
        {#if replaced[i]}
        {/if}
    {/each}
</main>