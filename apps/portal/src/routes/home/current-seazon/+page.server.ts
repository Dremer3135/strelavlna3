import { fail, redirect, error } from '@sveltejs/kit';
import { enhance } from '$app/forms';
import { createPocketbaseInstance } from '$lib/server/pocketbase';


export const load = async ({ locals }) => {
	try {
        let [contests, teams] = await Promise.all([
            locals.pb.collection('contests').getFullList({ filter: "public=true"}),
            locals.pb.collection('teams').getFullList({ filter: `teacher.id='${locals.pb.authStore.record?.id}'`})
        ]);

        return {
			contests: contests,
			teams: teams
		};
	} catch (err) {
		console.error('Error in probeditor load function:', err);
		throw error(500, 'Something went wrong while fetching data.');
	}
};
