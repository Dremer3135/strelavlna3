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

<style>
    @import url('https://fonts.googleapis.com/css2?family=Lexend:wght@100..900&display=swap');
    @font-face {
        font-family: "Fredoka";
        
        src: url('/fonts/Fredoka-VF.ttf') format('truetype');

        font-style: normal;
        font-weight: normal;
    }

    :root {
        --color-purple: #9500EB;
        --color-pink: #EB0072;
        --color-orange: #EB6E00;
        --color-yellow: #EBAD00;
        --color-blue: #3118ba;
        --color-lightblue:#4c33da;
        --color-gray: #353538;
    }




</style>
