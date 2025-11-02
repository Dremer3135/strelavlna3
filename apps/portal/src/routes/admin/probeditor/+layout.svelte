<script lang="ts">
    import { editableProbs } from "$lib/stores/probs";
    import type { EditableProb } from '$lib/types';

    let { children, data } = $props();

    editableProbs.set(Object.fromEntries(data.probs.map(prob => [prob.id, {prob, edit: {}}])));

    $effect(() => {
        editableProbs.set(Object.fromEntries(data.probs.map(prob => [prob.id, {prob, edit: {}}])));
    });

</script>

{@render children()}
