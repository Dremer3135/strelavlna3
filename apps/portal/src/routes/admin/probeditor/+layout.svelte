<script lang="ts">
    import { subscribeToConstants } from '$lib/stores/consts.js';
    import { subscribeToProbs } from '$lib/stores/probs';
    import { currentUser } from '$lib/stores/auth.js';
    import { pocketbase } from '$lib/pocketbase.js';

    let { children, data } = $props();

    $effect(() => {
        let unsubscribeProbs: (() => void) | undefined;
        let unsubscribeConstants: (() => void) | undefined;

        // This effect will re-run whenever the user logs in or out.
        if ($currentUser) {
            // User is authenticated, we can now subscribe.
            console.log('User is authenticated, setting up subscriptions...', pocketbase.authStore.model);

            const setupSubscriptions = async () => {
                try {
                    unsubscribeProbs = await subscribeToProbs(data.probs);
                    unsubscribeConstants = await subscribeToConstants(data.constants);
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
        };
    });
</script>

{@render children()}


