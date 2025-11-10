import PocketBase from 'pocketbase';
import type { TypedPocketBase } from '$lib/types/pocketbase-types';
import { POCKETBASE_URL } from '$env/static/private';

export const createPocketbaseInstance = () => {
    if (!POCKETBASE_URL) {
        throw new Error('The POCKETBASE_URL environment variable is not set. Please set it in your .env file.');
    }
    return new PocketBase(POCKETBASE_URL) as TypedPocketBase;
}
