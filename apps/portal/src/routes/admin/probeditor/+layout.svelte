<script lang="ts">
    import { onMount } from 'svelte';
    import { probs } from "$lib/stores/probs";
    let { children, data } = $props();

    $effect(() => {
        probs.set(data.probs);
    });

</script>

{@render children()}
