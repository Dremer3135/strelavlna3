<script lang="ts">
    import { editableConsts } from "$lib/stores/consts.js";
    import { editableProbs } from "$lib/stores/probs";
    import type { EditableProb } from '$lib/types';

    let { children, data } = $props();

    editableProbs.set(Object.fromEntries(data.probs.map(prob => [prob.id, {prob, edit: {}}])));
    $effect(() => {
        editableProbs.set(Object.fromEntries(data.probs.map(prob => [prob.id, {prob, edit: {}}])));
    });
    
    // editableConsts.set(Object.fromEntries(data.consts.map(prob => [const.id, {const, edit: {}}])));

</script>

{@render children()}
