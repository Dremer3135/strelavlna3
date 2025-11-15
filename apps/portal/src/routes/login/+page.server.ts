import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    default: async ({ request, locals, url }) => {
        const formData = await request.formData();
        const email = formData.get('email') as string;
        const password = formData.get('password') as string;
        const redirectTo = formData.get('redirectTo');
        const adminLogin = (formData.get('adminLogin') as string) === 'true';

        console.log("admiiiin:", adminLogin);
        try {
            
            await locals.pb.collection(adminLogin ? 'correctors' : 'teachers').authWithPassword(email, password);
        } catch (err) {
            console.error(err);
            return fail(400, { errorType: 'invalid_credentials' });
        }

        if (redirectTo) {
            throw redirect(303, redirectTo.toString());
        }

        throw redirect(303, '/');
    }
};

// We need to combine the load function and the actions into the same file.
export const load = ({ url, locals }) => {
    if (locals.user && (locals.user.collectionName === 'correctors' || locals.user.collectionName === 'teachers')) {
        throw redirect(303, '/');
    }
    const redirectTo = url.searchParams.get('redirectTo');
    const adminLogin: boolean = url.searchParams.get('adminLogin') === 'true';
    return { redirectTo, adminLogin };
};
