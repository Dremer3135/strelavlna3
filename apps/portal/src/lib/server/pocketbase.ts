import PocketBase from 'pocketbase';
import type { TypedPocketBase } from '$lib/pocketbase-types';
import { POCKETBASE_URL } from '$env/static/private';

export const createPocketbaseInstance = () => {
    return new PocketBase(POCKETBASE_URL) as TypedPocketBase;
}