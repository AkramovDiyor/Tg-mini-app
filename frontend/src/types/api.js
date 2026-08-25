/**
 * JSDoc-контракты API. Полный TS (.tsx) — отдельным шагом.
 * @typedef {'master' | 'client'} UserRole
 *
 * @typedef {Object} TelegramUser
 * @property {number} [telegram_id]
 * @property {string} [first_name]
 * @property {string} [last_name]
 * @property {string} [username]
 *
 * @typedef {Object} WorkHours
 * @property {boolean[]} [work_days]
 * @property {string} [start_time]
 * @property {string} [end_time]
 * @property {string} [lunch_start]
 * @property {string} [lunch_end]
 *
 * @typedef {Object} MasterSettings
 * @property {boolean} [auto_cancel]
 * @property {string} [cancel_hours]
 * @property {boolean} [offer_waitlist]
 *
 * @typedef {Object} Master
 * @property {number} id
 * @property {number} telegram_id
 * @property {string} name
 * @property {string} [bio]
 * @property {string} [address]
 * @property {string} invite_link
 * @property {WorkHours} [work_hours]
 * @property {MasterSettings} [settings]
 *
 * @typedef {Object} MasterPhoto
 * @property {number} id
 * @property {number} [master_id]
 * @property {string} url
 *
 * @typedef {Object} MasterInfo
 * @property {number} id
 * @property {string} name
 * @property {string} [bio]
 * @property {string} [address]
 * @property {string} [invite_link]
 * @property {number} [rating]
 * @property {MasterPhoto[]} [photos]
 *
 * @typedef {Object} Identity
 * @property {UserRole} role
 * @property {Master} [master]
 * @property {TelegramUser} [user]
 * @property {string} [invite_link]
 *
 * @typedef {Object} Service
 * @property {number} id
 * @property {number} [master_id]
 * @property {string} name
 * @property {number} duration_min
 * @property {number} price
 *
 * @typedef {Object} Slot
 * @property {string} start_time
 * @property {string} status
 *
 * @typedef {Object} BookSlotRequest
 * @property {string} start_time
 * @property {number} service_id
 * @property {string} name
 * @property {number} price
 * @property {string} [invite_link]
 * @property {string} [date]
 *
 * @typedef {Object} ClientBooking
 * @property {number} booking_id
 * @property {string} client_name
 * @property {string} service_name
 * @property {number} service_duration
 * @property {number} service_price
 * @property {string} master_name
 * @property {string} master_address
 * @property {string} [master_invite_link]
 * @property {string} start_time
 * @property {string} end_time
 * @property {string} status
 * @property {string} created_at
 *
 * @typedef {Object} TodaySchedule
 * @property {{ count: number, total: number }} stats
 * @property {Array<Object>} schedule
 *
 * @typedef {Object} WaitlistItem
 * @property {number} id
 * @property {string} [client_name]
 * @property {string} [desired_date]
 * @property {string} [created_at]
 *
 * @typedef {Object} UpdateSettingsRequest
 * @property {WorkHours} work_hours
 * @property {MasterSettings} settings
 */

export {}
