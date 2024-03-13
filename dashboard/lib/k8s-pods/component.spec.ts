import {expect} from '@playwright/test';
import chalk from "chalk";

export const testK8sPodsComponentStep = async ({page}) => {
    console.log(chalk.green('K8s Pods Component Step'));
    // TODO: authorize and write tests
    await page.goto('/');
    await expect(page).toHaveTitle(/Grafana/);
}
