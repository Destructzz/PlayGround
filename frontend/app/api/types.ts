export interface ApiResponseMeta {
  timestamp: string
  request_id?: string
}

export interface AuthUser {
  id: string
  email: string
  avatar_url: string
  name: string
  provider: string
}

export interface SessionResponse extends ApiResponseMeta {
  authenticated: boolean
  user: AuthUser | null
}

export interface CatalogZone {
  id: number
  name: string
  zone_type: string
  zone_tag_id: number
  capacity: number
  description: string
  is_active: boolean
  details_json: unknown
}

export interface CatalogService {
  id: number
  name: string
  duration: number
  price: string
  currency: string
  description: string
  details_json: unknown
}

export interface CatalogZoneWithServices extends CatalogZone {
  services: CatalogService[]
}

export interface HomeCatalogResponse extends ApiResponseMeta {
  gaming: CatalogZone[]
  lounge: CatalogZone[]
  event: CatalogZone[]
}

export interface LoungeCatalogResponse extends ApiResponseMeta {
  zones: CatalogZoneWithServices[]
}

export interface EventCatalogResponse extends ApiResponseMeta {
  zones: CatalogZoneWithServices[]
}

export interface GamingZoneTag {
  id: number
  name: string
}

export interface GamingConfiguration {
  id: number
  zone_tag_id: number
  specs_json: unknown
}

export interface GamingPlace {
  id: number
  label: string
  configuration_id?: number | null
  sort_order: number
  is_active: boolean
  specs?: unknown
}

export interface GamingZone extends CatalogZone {
  places: GamingPlace[]
  services: CatalogService[]
}

export interface GamingCatalogResponse extends ApiResponseMeta {
  zone_tags: GamingZoneTag[]
  zones: GamingZone[]
  configurations: GamingConfiguration[]
}

export interface ShiftUser {
  id: string
  full_name: string
  email: string
  avatar_url?: string | null
  phone?: string | null
  role: string
  is_active: boolean
}

export interface ShiftSchedule {
  id: number
  user_id: string
  zone_tag_id?: number | null
  start_time: string
  end_time: string
  note?: string | null
  created_at: string
  updated_at: string
  user: ShiftUser
}

export interface ShiftListResponse extends ApiResponseMeta {
  shifts: ShiftSchedule[]
}
