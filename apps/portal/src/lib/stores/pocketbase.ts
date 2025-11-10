import { writable } from 'svelte/store';
import type { TypedPocketBase } from '$lib/types/pocketbase-types';

export const pb = writable<TypedPocketBase | undefined | null>(undefined);