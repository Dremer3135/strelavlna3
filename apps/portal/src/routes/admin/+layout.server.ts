import { redirect } from '@sveltejs/kit';

export const load = async ({ locals, url }) => {
    if (!locals.user) {
        throw redirect(303, `/admin-login?redirectTo=${url.pathname}`);
    }

    if (locals.user.collectionName !== 'correctors') {
        throw redirect(303, '/');
    }
};