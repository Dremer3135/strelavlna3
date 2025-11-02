<script lang="ts">
    import Navbar from '$lib/components/admin/Navbar.svelte';
    let { data, children } = $props();
    import type { CorrectorsResponse } from '$lib/pocketbase-types';
</script>

<div class="admin-layout-container">
    <Navbar user={data.user as CorrectorsResponse} isAdmin={false}/>
    {@render children()}
</div>

<style>
    .admin-layout-container {
        display: flex;
        flex-direction: column;
        height: 100vh; /* Ensure the container takes full viewport height */
    }
</style>
