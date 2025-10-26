<script lang="ts">
    import Navbar from '$lib/components/admin/Navbar.svelte';
    let { data, children } = $props();
    import type { CorrectorsResponse } from '$lib/pocketbase-types';
</script>

<Navbar user={data.user as CorrectorsResponse} />

<main style="padding: 1rem;">
    {@render children()}
</main>
