import type { CorrectorsResponse, ProbsResponse, TeachersResponse } from './pocketbase-types';
import type { EditableProb } from './types';

type AppUser = CorrectorsResponse | TeachersResponse;

/**
 * Type Guard to check if a user is a Corrector.
 * @param user The user object to check.
 * @returns True if the user is a CorrectorsResponse, false otherwise.
 */
export const isCorrector = (user: AppUser | undefined): user is CorrectorsResponse => {
    return user?.collectionName === 'correctors';
}


export const isEdited = (eprob: EditableProb): boolean => {
    return Object.keys(eprob.edit).length > 0
}

export const getEditedState = (eprob: EditableProb): ProbsResponse => {
    return {
        ...eprob.prob,
        ...(eprob.edit as Partial<ProbsResponse>)
    };
};

export function filterRecord<T>( 
    record: Record<string, T>,
    predicate: (id: string, value: T) => boolean
): Record<string, T> {
    return Object.entries(record)
        .filter(([id, value]) => predicate(id, value as T))
        .reduce((acc, [id, value]) => {
            acc[id] = value as T;
            return acc;
        }, {} as Record<string, T>);
}
