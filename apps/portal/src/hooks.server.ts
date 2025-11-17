import { createPocketbaseInstance } from '$lib/server/pocketbase';
import { redirect } from '@sveltejs/kit';
import type { AppUser } from './app';
// import { isCorrector } from '$lib/utils';

export const handle = async ({ event, resolve }) => {
	event.locals.pb = createPocketbaseInstance();
	event.locals.pb.authStore.loadFromCookie(event.request.headers.get('cookie') || '');

	event.locals.pb.beforeSend = (url, options) => {
		options.headers = Object.assign({}, options.headers, {
			'Origin': event.url.origin,
		});

		return { url, options };
	};
	
  if (event.url.pathname === '/qr') {
    throw redirect(302, '/');
  }

  if (event.url.pathname === '/logins') {
    throw redirect(302, '/home/register');
  }
	
	if (event.url.pathname === '/logout' && event.request.method === 'POST') {
		const userIsCorrector = event.locals.pb.authStore.model?.collectionName === 'correctors';

		// Clear the server-side auth store
		event.locals.pb.authStore.clear();

		// Clear the browser's cookie
		event.cookies.delete('pb_auth', { path: '/' });

		// if (userIsCorrector) {
		throw redirect(303, '/');
		// } else {
		// 	throw redirect(303, '/login');
		// }
	}

	try {
		// get an up-to-date auth store state by verifying and refreshing the loaded auth model (if any)
		if (event.locals.pb.authStore.isValid && event.locals.pb.authStore.model) {
			await event.locals.pb.collection(event.locals.pb.authStore.model.collectionName).authRefresh();
			event.locals.user = event.locals.pb.authStore.model as AppUser;
		}
	} catch (_) {
		// clear the auth store on failed refresh
		event.locals.pb.authStore.clear();
	}

	// event.locals.isAdminLogin = false;
	if (event.url.pathname === '/login/admin') {
		throw redirect(307, '/login?adminLogin=true');
	} else if (event.url.pathname === '/login' && !event.url.searchParams.has('adminLogin')) {
		event.url.searchParams.set('adminLogin', 'false');
	}

	const response = await resolve(event);

	// send back the default 'pb_auth' cookie to the client with the latest store state
	response.headers.append('set-cookie', event.locals.pb.authStore.exportToCookie());

	return response;
};
