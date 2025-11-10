import PocketBase from 'pocketbase';
import { PUBLIC_POCKETBASE_URL } from '$env/static/public';
import type { TypedPocketBase } from '$lib/types/pocketbase-types';
import { currentUser } from '$lib/stores/auth';

export const pocketbase = new PocketBase(PUBLIC_POCKETBASE_URL) as TypedPocketBase;

// Set the initial value of the authStore
currentUser.set(pocketbase.authStore.model);

pocketbase.authStore.onChange(() => {
    currentUser.set(pocketbase.authStore.model);
});