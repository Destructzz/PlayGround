export type MockUiState = 'default' | 'loading' | 'empty' | 'error' | 'success'

export interface LoungeZone {
  id: string
  name: string
  mood: string
  capacity: number
  remaining: number
  description: string
  perks: string[]
  accent: string
  premium?: boolean
}

export interface EventItem {
  id: string
  title: string
  dateLabel: string
  timeLabel: string
  format: string
  remaining: number
  capacity: number
  description: string
  accent: string
  speakers?: string[]
  soldOut?: boolean
}

export interface SelectOption {
  id: string
  label: string
  disabled?: boolean
}

export const mockStateLabels: Record<MockUiState, { title: string, description: string }> = {
  default: {
    title: 'PlayGround Experience',
    description: 'Полноценные mock-сценарии без привязки к backend.'
  },
  loading: {
    title: 'Подготавливаем интерфейс',
    description: 'Показываем deterministic loading-state для проверки UX до подключения API.'
  },
  empty: {
    title: 'Пока ничего не запланировано',
    description: 'Пустое состояние должно оставаться информативным и не ломать маршрут.'
  },
  error: {
    title: 'Нужна повторная попытка',
    description: 'Ошибка здесь симулируется локально, чтобы заранее отработать edge-case.'
  },
  success: {
    title: 'Сценарий уже завершён',
    description: 'Предпросмотр success-state доступен напрямую через query-параметр.'
  }
}

export const loungeZones: LoungeZone[] = [
  {
    id: 'aurora',
    name: 'Aurora Corner',
    mood: 'Quiet social',
    capacity: 6,
    remaining: 4,
    description: 'Полумягкая lounge-зона с приватным светом, низкими столами и комфортной посадкой для долгих вечерних сессий.',
    perks: ['Signature drinks', 'Ambient lighting', 'PS5 side station'],
    accent: '#22d3ee'
  },
  {
    id: 'ember',
    name: 'Ember Booth',
    mood: 'Premium private',
    capacity: 8,
    remaining: 2,
    description: 'Закрытая секция с более плотной посадкой, тёплым оранжевым контуром и приоритетным сервисом для небольших компаний.',
    perks: ['Private audio', 'Priority service', 'Reserved host support'],
    accent: '#fb923c',
    premium: true
  },
  {
    id: 'pulse',
    name: 'Pulse Hall',
    mood: 'Open energy',
    capacity: 12,
    remaining: 0,
    description: 'Открытая зона у главного потока гостей: лучший обзор на ивенты, но на популярные вечера места заканчиваются раньше всего.',
    perks: ['Main stage view', 'Fast refill lane', 'Community seating'],
    accent: '#a3e635'
  }
]

export const loungeTimeSlots: SelectOption[] = [
  { id: '18:00', label: '18:00' },
  { id: '19:00', label: '19:00' },
  { id: '20:00', label: '20:00', disabled: true },
  { id: '21:00', label: '21:00' },
  { id: '22:00', label: '22:00' }
]

export const loungePartyOptions: SelectOption[] = [
  { id: '2', label: '2 гостя' },
  { id: '4', label: '4 гостя' },
  { id: '6', label: '6 гостей' },
  { id: '8', label: '8 гостей', disabled: true }
]

export const eventItems: EventItem[] = [
  {
    id: 'night-bracket',
    title: 'Night Bracket Qualifier',
    dateLabel: 'Сегодня',
    timeLabel: '19:30',
    format: 'Tournament',
    remaining: 5,
    capacity: 32,
    description: 'Вечерняя турнирная сетка для команд и соло-участников с живым комментированием и LED-сценой.',
    accent: '#22d3ee',
    speakers: ['Host: Raven', 'Caster: Miko']
  },
  {
    id: 'creator-lab',
    title: 'Creator Lab Meetup',
    dateLabel: 'Завтра',
    timeLabel: '18:00',
    format: 'Meetup',
    remaining: 12,
    capacity: 40,
    description: 'Интенсив для локальных контент-криэйторов, брендов и организаторов community-ивентов.',
    accent: '#f472b6',
    speakers: ['Lera Bloom', 'Den Vox']
  },
  {
    id: 'full-house',
    title: 'Retro LAN Story Night',
    dateLabel: 'Пятница',
    timeLabel: '21:00',
    format: 'Community event',
    remaining: 0,
    capacity: 24,
    description: 'Ламповый late-night формат с ретро-сетапами и curated плейлистом. Все места уже заняты.',
    accent: '#f59e0b',
    soldOut: true,
    speakers: ['Guest host: Tair']
  }
]

export const eventGuestOptions: SelectOption[] = [
  { id: 'solo', label: 'Соло' },
  { id: 'duo', label: 'Дуо' },
  { id: 'squad', label: 'С компанией' }
]
