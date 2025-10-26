// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
import type { TypedPocketBase } from "$lib/pocketbase-types";
import type { CorrectorsResponse, TeachersResponse } from "$lib/pocketbase-types";

export type AppUser = CorrectorsResponse | TeachersResponse;

declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			pb: TypedPocketBase,
			user: AppUser | undefined
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
