<script lang="ts">
    import Navbar from '$lib/components/admin/Navbar.svelte';
<<<<<<< HEAD
    let { data, children } = $props();
    import type { CorrectorsResponse } from '$lib/types/pocketbase-types';
=======
    import type { CorrectorsResponse } from '$lib/pocketbase-types';
    import { subscribeToConstants } from '$lib/stores/consts.js';
    import { subscribeToContests } from '$lib/stores/contests.js';
    import { subscribeToProbs } from '$lib/stores/probs';
    import { correctors } from '$lib/stores/correctors.js';
    import { currentUser } from '$lib/stores/auth.js';
    import { pocketbase } from '$lib/pocketbase.js';

    let { data, children } = $props();

    $effect(() => {
        correctors.set(Object.fromEntries(data.correctors.map((corrector) => [corrector.id, corrector])));
        let unsubscribeProbs: (() => void) | undefined;
        let unsubscribeConstants: (() => void) | undefined;
        let unsubscribeContests: (() => void) | undefined;

        // This effect will re-run whenever the user logs in or out.
        if ($currentUser) {
            // User is authenticated, we can now subscribe.
            console.log('User is authenticated, setting up subscriptions...', pocketbase.authStore.model);

            const setupSubscriptions = async () => {
                try {
                    unsubscribeProbs = await subscribeToProbs(data.probs);
                    unsubscribeConstants = await subscribeToConstants(data.constants);
                    unsubscribeContests = await subscribeToContests(data.contests);
                } catch (error) {
                    console.error("Failed to subscribe:", error);
                }
            };

            setupSubscriptions();
        }

        // The return function from $effect is the cleanup function.
        return () => {
            if (unsubscribeProbs) {
                console.log('Unsubscribing from probs.');
                unsubscribeProbs();
            }
            if (unsubscribeConstants) {
                console.log('Unsubscribing from constants.');
                unsubscribeConstants();
            }
            if (unsubscribeContests) {
                console.log('Unsubscribing from contests.');
                unsubscribeContests();
            }
        };
    });


>>>>>>> 611dff672c0713e762451f75c117fd7a93ae168c
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
