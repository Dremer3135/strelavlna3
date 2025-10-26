<script lang="ts">
    import type { CorrectorsResponse } from "$lib/pocketbase-types";

    let { user }: { user: CorrectorsResponse | undefined } = $props();
</script>

<header style="display: flex; justify-content: space-between; align-items: center; padding: 1rem; border-bottom: 1px solid #ccc; background-color: #f9f9f9;">
    <div>
        <strong>Admin Panel</strong>
        {#if user}
            <span style="margin-left: 1rem; font-size: 0.9rem; color: #555;">Signed in as {user.email}</span>
        {/if}
    </div>
    {#if user}
        <form method="POST" action="/logout">
            <button type="submit">Logout</button>
        </form>
    {/if}
</header>
