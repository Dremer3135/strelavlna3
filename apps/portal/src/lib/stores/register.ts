import { writable, readable } from 'svelte/store';
import type { Contest, Team } from '$lib/types/register';
import { pocketbase } from '$lib/pocketbase';
import type { RecordModel } from 'pocketbase';
import { currentUser } from '$lib/stores/auth';

// Contests Store
export const contestsStore = readable<Contest[]>([], (set) => {
    pocketbase.collection('contests').getFullList<Contest>({ sort: '-created', filter: 'public = true' })
        .then(contests => {
            set(contests);
        })
        .catch(err => {
            console.error('Error fetching contests:', err);
            set([]);
        });
});

// Teams Store
export const teamsStore = writable<Team[]>([]);

export async function fetchTeams(teacherId: string) {
    try {
        const teams = await pocketbase.collection('teams').getFullList<Team>({
            filter: `teacher = "${teacherId}"`,
        });
        teamsStore.set(teams);
    } catch (err) {
        console.error('Error fetching teams:', err);
    }
}
