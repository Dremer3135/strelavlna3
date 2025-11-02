import { writable } from 'svelte/store';
// import type { AppUser } from '$lib/pocketbase-types';
import type { AppUser } from '../../app';

export const currentUser = writable<AppUser | undefined | null>(undefined);