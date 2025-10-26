<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
    import { page } from '$app/stores';
    import { currentUser } from '$lib/stores/auth';

	let { children } = $props();

    $effect(() => {
        currentUser.set($page.data.user);
    });
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children?.()}
