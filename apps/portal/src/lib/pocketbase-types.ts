/**
* This file was @generated using pocketbase-typegen
*/

import type PocketBase from 'pocketbase'
import type { RecordService } from 'pocketbase'

export enum Collections {
	Authorigins = "_authOrigins",
	Externalauths = "_externalAuths",
	Mfas = "_mfas",
	Otps = "_otps",
	Superusers = "_superusers",
	Consts = "consts",
	Contests = "contests",
	Correctors = "correctors",
	Probs = "probs",
	Skoly = "skoly",
	Teachers = "teachers",
	Teams = "teams",
	Texts = "texts",
}

// Alias types for improved usability
export type IsoDateString = string
export type RecordIdString = string
export type HTMLString = string

type ExpandType<T> = unknown extends T
	? T extends unknown
		? { expand?: unknown }
		: { expand: T }
	: { expand: T }

// System fields
export type BaseSystemFields<T = unknown> = {
	id: RecordIdString
	collectionId: string
	collectionName: Collections
} & ExpandType<T>

export type AuthSystemFields<T = unknown> = {
	email: string
	emailVisibility: boolean
	username: string
	verified: boolean
} & BaseSystemFields<T>

// Record types for each collection

export type AuthoriginsRecord = {
	collectionRef: string
	created?: IsoDateString
	fingerprint: string
	id: string
	recordRef: string
	updated?: IsoDateString
}

export type ExternalauthsRecord = {
	collectionRef: string
	created?: IsoDateString
	id: string
	provider: string
	providerId: string
	recordRef: string
	updated?: IsoDateString
}

export type MfasRecord = {
	collectionRef: string
	created?: IsoDateString
	id: string
	method: string
	recordRef: string
	updated?: IsoDateString
}

export type OtpsRecord = {
	collectionRef: string
	created?: IsoDateString
	id: string
	password: string
	recordRef: string
	sentTo?: string
	updated?: IsoDateString
}

export type SuperusersRecord = {
	created?: IsoDateString
	email: string
	emailVisibility?: boolean
	id: string
	password: string
	tokenKey: string
	updated?: IsoDateString
	verified?: boolean
}

export type ConstsRecord = {
	created?: IsoDateString
	desc?: string
	group?: string
	id: string
	name?: string
	symbol?: string
	unit?: string
	updated?: IsoDateString
	value?: number
}

export enum ContestsTypeOptions {
	"math" = "math",
	"physics" = "physics",
}
export type ContestsRecord = {
	created?: IsoDateString
	id: string
	name?: string
	onSiteEnd?: IsoDateString
	onSiteStart?: IsoDateString
	onlineEnd?: IsoDateString
	onlineStart?: IsoDateString
	public?: boolean
	type?: ContestsTypeOptions
	updated?: IsoDateString
}

export type CorrectorsRecord = {
	admin?: boolean
	created?: IsoDateString
	email?: string
	emailVisibility?: boolean
	id: string
	password: string
	socials?: string
	tokenKey: string
	updated?: IsoDateString
	username: string
	verified?: boolean
}

export enum ProbsDiffOptions {
	"A" = "A",
	"B" = "B",
	"C" = "C",
}
export type ProbsRecord = {
	answer?: string
	author?: RecordIdString
	auto?: boolean
	code?: string
	contests?: RecordIdString[]
	created?: IsoDateString
	diff?: ProbsDiffOptions
	id: string
	infinite?: boolean
	name?: string
	queue?: RecordIdString[]
	text?: string
	updated?: IsoDateString
}

export type SkolyRecord = {
	c_obce?: string
	c_or?: number
	c_p?: number
	created?: IsoDateString
	datum_zahajeni_cinnosti?: IsoDateString
	email_1?: string
	email_2?: string
	fax?: string
	ico?: number
	id: string
	izo?: number
	kapacita?: number
	kod_ruian?: number
	kraj?: string
	misto?: string
	okres?: string
	orp?: number
	orp_nazev?: string
	plny_nazev?: string
	psc?: number
	red_izo?: number
	reditel?: string
	spam?: boolean
	spravni_urad?: string
	telefon?: string
	ulice?: string
	updated?: IsoDateString
	uzemi?: string
	www?: string
	zkraceny_nazev?: string
	zrizovatel?: number
}

export type TeachersRecord = {
	created?: IsoDateString
	email?: string
	emailVisibility?: boolean
	id: string
	password: string
	skola?: RecordIdString
	tokenKey: string
	updated?: IsoDateString
	username: string
	verified?: boolean
}

export type TeamsRecord = {
	contest?: RecordIdString
	created?: IsoDateString
	id: string
	name?: string
	player1email?: string
	player1name?: string
	player2email?: string
	player2name?: string
	player3email?: string
	player3name?: string
	player4email?: string
	player4name?: string
	player5email?: string
	player5name?: string
	teacher?: RecordIdString
	updated?: IsoDateString
}

export type TextsRecord = {
	created?: IsoDateString
	id: string
	name?: string
	text?: HTMLString
	updated?: IsoDateString
}

// Response types include system fields and match responses from the PocketBase API
export type AuthoriginsResponse<Texpand = unknown> = Required<AuthoriginsRecord> & BaseSystemFields<Texpand>
export type ExternalauthsResponse<Texpand = unknown> = Required<ExternalauthsRecord> & BaseSystemFields<Texpand>
export type MfasResponse<Texpand = unknown> = Required<MfasRecord> & BaseSystemFields<Texpand>
export type OtpsResponse<Texpand = unknown> = Required<OtpsRecord> & BaseSystemFields<Texpand>
export type SuperusersResponse<Texpand = unknown> = Required<SuperusersRecord> & AuthSystemFields<Texpand>
export type ConstsResponse<Texpand = unknown> = Required<ConstsRecord> & BaseSystemFields<Texpand>
export type ContestsResponse<Texpand = unknown> = Required<ContestsRecord> & BaseSystemFields<Texpand>
export type CorrectorsResponse<Texpand = unknown> = Required<CorrectorsRecord> & AuthSystemFields<Texpand>
export type ProbsResponse<Texpand = unknown> = Required<ProbsRecord> & BaseSystemFields<Texpand>
export type SkolyResponse<Texpand = unknown> = Required<SkolyRecord> & BaseSystemFields<Texpand>
export type TeachersResponse<Texpand = unknown> = Required<TeachersRecord> & AuthSystemFields<Texpand>
export type TeamsResponse<Texpand = unknown> = Required<TeamsRecord> & BaseSystemFields<Texpand>
export type TextsResponse<Texpand = unknown> = Required<TextsRecord> & BaseSystemFields<Texpand>

// Types containing all Records and Responses, useful for creating typing helper functions

export type CollectionRecords = {
	_authOrigins: AuthoriginsRecord
	_externalAuths: ExternalauthsRecord
	_mfas: MfasRecord
	_otps: OtpsRecord
	_superusers: SuperusersRecord
	consts: ConstsRecord
	contests: ContestsRecord
	correctors: CorrectorsRecord
	probs: ProbsRecord
	skoly: SkolyRecord
	teachers: TeachersRecord
	teams: TeamsRecord
	texts: TextsRecord
}

export type CollectionResponses = {
	_authOrigins: AuthoriginsResponse
	_externalAuths: ExternalauthsResponse
	_mfas: MfasResponse
	_otps: OtpsResponse
	_superusers: SuperusersResponse
	consts: ConstsResponse
	contests: ContestsResponse
	correctors: CorrectorsResponse
	probs: ProbsResponse
	skoly: SkolyResponse
	teachers: TeachersResponse
	teams: TeamsResponse
	texts: TextsResponse
}

// Type for usage with type asserted PocketBase instance
// https://github.com/pocketbase/js-sdk#specify-typescript-definitions

export type TypedPocketBase = PocketBase & {
	collection(idOrName: '_authOrigins'): RecordService<AuthoriginsResponse>
	collection(idOrName: '_externalAuths'): RecordService<ExternalauthsResponse>
	collection(idOrName: '_mfas'): RecordService<MfasResponse>
	collection(idOrName: '_otps'): RecordService<OtpsResponse>
	collection(idOrName: '_superusers'): RecordService<SuperusersResponse>
	collection(idOrName: 'consts'): RecordService<ConstsResponse>
	collection(idOrName: 'contests'): RecordService<ContestsResponse>
	collection(idOrName: 'correctors'): RecordService<CorrectorsResponse>
	collection(idOrName: 'probs'): RecordService<ProbsResponse>
	collection(idOrName: 'skoly'): RecordService<SkolyResponse>
	collection(idOrName: 'teachers'): RecordService<TeachersResponse>
	collection(idOrName: 'teams'): RecordService<TeamsResponse>
	collection(idOrName: 'texts'): RecordService<TextsResponse>
}
