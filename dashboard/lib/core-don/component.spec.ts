import {expect} from '@playwright/test';
import chalk from 'chalk';


export const testCoreDonComponentStep = async ({page}) => {
    console.log(chalk.green('Core DON Component Step'));
    // TODO: authorize and write tests
    await page.goto('/');
    await expect(page).toHaveTitle(/Grafana/);
}