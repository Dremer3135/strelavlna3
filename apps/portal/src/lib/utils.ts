import type { CorrectorsResponse, TeachersResponse } from './pocketbase-types';

type AppUser = CorrectorsResponse | TeachersResponse;

/**
 * Type Guard to check if a user is a Corrector.
 * @param user The user object to check.
 * @returns True if the user is a CorrectorsResponse, false otherwise.
 */
export const isCorrector = (user: AppUser | undefined): user is CorrectorsResponse => {
    return user?.collectionName === 'correctors';
}
