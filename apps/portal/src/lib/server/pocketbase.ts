import PocketBase from 'pocketbase';
import type { TypedPocketBase } from '$lib/pocketbase-types';

export const createPocketbaseInstance = () => {
    return new PocketBase('https://strela-vlna.gchd.cz') as TypedPocketBase;
}
