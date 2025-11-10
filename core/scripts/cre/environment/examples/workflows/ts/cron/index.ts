import { cre, Runner, type Runtime } from '@chainlink/cre-sdk'

type Config = {
    schedule: string
}

const onCronTrigger = (runtime: Runtime<Config>): string => {
    runtime.log('Hello world! Workflow triggered.')
    return 'Hello world!'
}

const initWorkflow = (config: Config) => {
    const cron = new cre.capabilities.CronCapability()

    return [cre.handler(cron.trigger({ schedule: config.schedule }), onCronTrigger)]
}

export async function main() {
    const runner = await Runner.newRunner<Config>()
    await runner.run(initWorkflow)
}

main()
