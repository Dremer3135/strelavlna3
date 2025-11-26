import type { ConstantsResponse, ContestsResponse, CorrectorsResponse, ProbsResponse, TeachersResponse } from './pocketbase-types';
import type { EditableConstant, EditableContest, EditableProb } from './types';

type AppUser = CorrectorsResponse | TeachersResponse;

/**
 * Type Guard to check if a user is a Corrector.
 * @param user The user object to check.
 * @returns True if the user is a CorrectorsResponse, false otherwise.
 */
export const isCorrector = (user: AppUser | undefined): user is CorrectorsResponse => {
    return user?.collectionName === 'correctors';
}


export const isProbEdited = (eprob: EditableProb): boolean => {
    return Object.keys(eprob.edit).length > 0
}

export const isConstantEdited = (econstant: EditableConstant): boolean => {
    return Object.keys(econstant.edit).length > 0
}

export const getProbEditedState = (eprob: EditableProb): ProbsResponse => {
    return {
        ...eprob.prob,
        ...(eprob.edit as Partial<ProbsResponse>)
    };
};

export const getConstantEditedState = (econstant: EditableConstant): ConstantsResponse => {
    return {
        ...econstant.constant,
        ...(econstant.edit as Partial<ConstantsResponse>)
    };
};

export const getContestEditedState = (econstant: EditableContest): ContestsResponse => {
    return {
        ...econstant.contest,
        ...(econstant.edit as Partial<ContestsResponse>)
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

export function removeDiacritics(text: string): string {
   return text
    .normalize('NFD') // Decompose combined characters into base character + diacritical mark
    .replace(/[\u0300-\u036f]/g, '') // Remove diacritical marks (Unicode range)
    .toLowerCase(); // Convert the entire string to lowercase
   }