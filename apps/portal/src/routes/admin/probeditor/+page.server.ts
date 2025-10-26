import { error } from '@sveltejs/kit';

export const load = async ({ locals }) => {
	try {
		const probs = await locals.pb.collection('probs').getFullList();
		const consts = await locals.pb.collection('consts').getFullList();

		return {
			probs: probs,
			consts: consts
		};
	} catch (err) {
		console.error('Error in probeditor load function:', err);
		throw error(500, 'Something went wrong while fetching data.');
	}
};