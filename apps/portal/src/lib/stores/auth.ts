import { writable } from 'svelte/store';
import type { UsersResponse } from '$lib/pocketbase-types';

export const currentUser = writable<UsersResponse | undefined | null>(undefined);