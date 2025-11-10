import { fail, redirect, error } from '@sveltejs/kit';
import { enhance } from '$app/forms';
import { createPocketbaseInstance } from '$lib/server/pocketbase';

import { POCKETBASE_SUPERUSER_EMAIL, POCKETBASE_SUPERUSER_PASSWORD } from '$env/static/private';

export const load = async ({ locals, url }) => {
    if (!locals.user) {
        throw redirect(303, `/admin-login?redirectTo=${url.pathname}`);
    }

    if (locals.user.collectionName !== 'correctors') {
        throw redirect(303, '/');
    }

    try {
		const [probs, constants, contests, correctors] = await Promise.all([
			locals.pb.collection('probs').getFullList(200, { requestKey: null }),
			locals.pb.collection('constants').getFullList(200, { requestKey: null }),
			locals.pb.collection('contests').getFullList(200, { requestKey: null }),
			locals.pb.collection('correctors').getFullList(200, { requestKey: null })
		]);

		return {
			probs: probs,
			constants: constants,
			contests: contests,
			correctors: correctors
		};
	} catch (err) {
		console.error('Error in probeditor load function:', err);
		throw error(500, 'Something went wrong while fetching data.');
	}

};