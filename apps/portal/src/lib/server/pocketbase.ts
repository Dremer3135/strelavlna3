import PocketBase from 'pocketbase';
import type { TypedPocketBase } from '$lib/pocketbase-types';

export const createPocketbaseInstance = () => {
    return new PocketBase('http://127.0.0.1:8090') as TypedPocketBase;
}