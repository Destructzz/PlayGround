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
  role: string
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

export type BookingStatus = 'created' | 'confirmed' | 'canceled' | 'completed'

export interface BookingRecord {
  id: number
  user_id: string
  zone_id: number
  service_id: number
  place_id?: number | null
  start_time: string
  end_time: string
  participants: number
  total_price: string
  status: BookingStatus
  contact_name: string
  contact_email: string
  contact_phone: string
  details_json: unknown
  created_at: string
  updated_at: string
}

export interface BookingResponse extends ApiResponseMeta {
  booking: BookingRecord
}

export interface BookingListResponse extends ApiResponseMeta {
  bookings: BookingRecord[]
}

export interface CreateBookingPayload {
  zone_id: number
  service_id: number
  place_id?: number
  start_time: string
  end_time: string
  participants: number
  status: BookingStatus
  contact_name: string
  contact_email: string
  contact_phone?: string
  details_json?: string
}

export interface PatchBookingPayload {
  zone_id?: number
  service_id?: number
  start_time?: string
  end_time?: string
  participants?: number
  status?: BookingStatus
  contact_name?: string
  contact_email?: string
  contact_phone?: string
}

export interface GamingAvailabilityBooking {
  booking_id: number
  place_id: number
  start_time: string
  end_time: string
  status: BookingStatus
}

export interface GamingAvailabilityResponse extends ApiResponseMeta {
  zone_id: number
  date: string
  bookings: GamingAvailabilityBooking[]
}

export interface LoungeSlot {
  hour: number
  label: string
  booked_participants: number
  remaining: number
  available: boolean
}

export interface LoungeAvailabilityResponse extends ApiResponseMeta {
  zone_id: number
  capacity: number
  date: string
  slots: LoungeSlot[]
}

export interface AdminService {
  id: number
  name: string
  zone_id: number
  duration: number
  price: string
  currency: string
  description: unknown
  is_active: boolean
  details_json: unknown
  created_at: string
  updated_at: string
}

export interface AdminServiceListResponse extends ApiResponseMeta {
  services: AdminService[]
}

export interface AdminZoneTag {
  id: number
  name: string
  created_at: string
  updated_at: string
}

export interface AdminZone {
  id: number
  name: string
  zone_type: string
  zone_tag_id: number
  capacity: number
  description: unknown
  is_active: boolean
  details_json: unknown
  created_at: string
  updated_at: string
}

export interface AdminPlace {
  id: number
  zone_id: number
  label: string
  configuration_id: unknown
  sort_order: number
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface AdminZoneTagListResponse extends ApiResponseMeta {
  zone_tags: AdminZoneTag[]
}

export interface AdminZoneListResponse extends ApiResponseMeta {
  zones: AdminZone[]
}

export interface AdminPlaceListResponse extends ApiResponseMeta {
  places: AdminPlace[]
}

export interface AdminZoneTagResponse extends ApiResponseMeta {
  zone_tag: AdminZoneTag
}

export interface AdminZoneResponse extends ApiResponseMeta {
  zone: AdminZone
}

export interface AdminPlaceResponse extends ApiResponseMeta {
  place: AdminPlace
}

export interface AdminBookingListResponse extends ApiResponseMeta {
  bookings: BookingRecord[]
}

export interface AdminBookingResponse extends ApiResponseMeta {
  booking: BookingRecord
}

export interface CreateZoneTagPayload {
  name: string
}

export interface PatchZoneTagPayload {
  name?: string
}

export interface CreateZonePayload {
  name: string
  type: string
  zone_tag_id: number
  capacity: number
  description: string
  is_active: boolean
  details_json: string
}

export interface PatchZonePayload {
  name?: string
  type?: string
  zone_tag_id?: number
  capacity?: number
  description?: string
  is_active?: boolean
  details_json?: string
}

export interface CreatePlacePayload {
  zone_id: number
  label: string
  configuration_id?: number
  sort_order?: number
  is_active?: boolean
}

export interface PatchPlacePayload {
  label?: string
  configuration_id?: number
  sort_order?: number
  is_active?: boolean
}

export interface CreateShiftPayload {
  start_time: string
  end_time: string
  zone_tag_id?: number
  note?: string
}

export interface PatchShiftPayload {
  start_time?: string
  end_time?: string
  zone_tag_id?: number
  note?: string
}

export interface CreateServicePayload {
  name: string
  zone_id: number
  duration: number
  price: number
  currency?: string
  description?: string
  is_active?: boolean
  details_json?: string
}

export interface PatchServicePayload {
  name?: string
  zone_id?: number
  duration?: number
  price?: number
  currency?: string
  description?: string
  is_active?: boolean
  details_json?: string
}

export interface AdminServiceResponse extends ApiResponseMeta {
  service: AdminService
}

