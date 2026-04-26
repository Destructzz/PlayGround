import { expect, test } from '@playwright/test'

test('homepage routes into lounge and event flows', async ({ page }) => {
  await page.goto('/')

  await page.getByTestId('home-lounge-cta').click()
  await expect(page).toHaveURL(/\/lounge$/)
  await expect(page.getByTestId('lounge-selected-zone')).toBeVisible()

  await page.goto('/')
  await page.getByTestId('home-event-cta').click()
  await expect(page).toHaveURL(/\/event$/)
  await expect(page.getByTestId('event-selected-card')).toBeVisible()
})

test('lounge flow enforces validation and completes happy path', async ({ page }) => {
  await page.goto('/lounge')
  await expect(page.getByTestId('lounge-client-ready')).toBeAttached()

  await page.getByTestId('lounge-submit').click()
  await expect(page.getByTestId('lounge-validation-message')).toBeVisible()

  await page.getByTestId('lounge-zone-aurora').click()
  await page.getByTestId('lounge-guest-name').fill('Alex')
  await page.getByTestId('lounge-contact').fill('@alex')
  await page.getByTestId('lounge-party-4').click()
  await page.getByTestId('lounge-slot-19:00').click()
  await page.getByTestId('lounge-submit').click()

  await expect(page.getByTestId('lounge-success-message')).toBeVisible()
})

test('event flow enforces validation and completes happy path', async ({ page }) => {
  await page.goto('/event')
  await expect(page.getByTestId('event-client-ready')).toBeAttached()

  await page.getByTestId('event-submit').click()
  await expect(page.getByTestId('event-validation-message')).toBeVisible()

  await page.getByTestId('event-card-night-bracket').click()
  await page.getByTestId('event-attendee-name').fill('Anna')
  await page.getByTestId('event-attendee-email').fill('anna@example.com')
  await page.getByTestId('event-attendance-solo').click()
  await page.getByTestId('event-submit').click()

  await expect(page.getByTestId('event-success-message')).toBeVisible()
})

test('error-state fixtures expose reachable submit failures for lounge and event', async ({ page }) => {
  await page.goto('/lounge?state=error')
  await expect(page.getByTestId('lounge-client-ready')).toBeAttached()
  await expect(page.getByTestId('lounge-error-mode')).toBeVisible()
  await page.getByTestId('lounge-zone-aurora').click()
  await page.getByTestId('lounge-guest-name').fill('Alex')
  await page.getByTestId('lounge-contact').fill('@alex')
  await page.getByTestId('lounge-party-4').click()
  await page.getByTestId('lounge-slot-19:00').click()
  await page.getByTestId('lounge-submit').click()
  await expect(page.getByTestId('lounge-submit-error')).toBeVisible()

  await page.goto('/event?state=error')
  await expect(page.getByTestId('event-client-ready')).toBeAttached()
  await expect(page.getByTestId('event-error-mode')).toBeVisible()
  await page.getByTestId('event-card-night-bracket').click()
  await page.getByTestId('event-attendee-name').fill('Anna')
  await page.getByTestId('event-attendee-email').fill('anna@example.com')
  await page.getByTestId('event-attendance-solo').click()
  await page.getByTestId('event-submit').click()
  await expect(page.getByTestId('event-submit-error')).toBeVisible()
})
