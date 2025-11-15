import { writable, readable } from 'svelte/store';
import type { ContestsResponse, TeamsResponse } from '$lib/pocketbase-types';
import { pocketbase } from '$lib/pocketbase';
import type { RecordModel } from 'pocketbase';

export const myTeamsStore = writable<TeamsResponse[]>([]);

export async function fetchTeams(teacherId: string) {
    try {
        const teams = await pocketbase.collection('teams').getFullList<TeamsResponse>({
            filter: `teacher = "${teacherId}"`,
        });
        myTeamsStore.set(teams);
    } catch (err) {
        console.error('Error fetching teams:', err);
    }
}
