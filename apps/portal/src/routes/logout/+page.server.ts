import { redirect } from '@sveltejs/kit';

export const load = async ({locals}) => {
    // const userIsCorrector = locals.pb.authStore.model?.collectionName === 'correctors';
    locals.pb.authStore.clear();
    throw redirect(303, '/');
}