<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import { currentUser } from '$lib/stores/auth';
	import { pb } from '$lib/pocketbase'; // Correctly import the client-side instance

	let { children, data } = $props();

	$effect(() => {
		// This syncs the client-side auth store with the cookie from the server
		pb.authStore.loadFromCookie(document.cookie || '');

		// This sets the Svelte store with the user data passed from the server
		currentUser.set(data.user);

		// Optional: You can add a listener to keep the Svelte store updated
		// if the auth state changes for any reason (e.g., manual login/logout).
		const removeListener = pb.authStore.onChange(() => {
			currentUser.set(pb.authStore.model as typeof data.user);
		});

		// Cleanup the listener when the component is destroyed
		return () => {
			removeListener();
		};
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

    :global(*, *::before, *::after) {
        box-sizing: border-box;
    }




</style>
