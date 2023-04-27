import { newEasiGen } from './index'
import userProfiles from '../cmd/easi-gen/internal/gen/UserProfiles.cdc?raw'

const main = async() => {
    const gen = await newEasiGen()
    const output = document.createElement('div')
    document.body.append(output)
    output.innerHTML = `<pre>${gen(userProfiles)}</pre>`
}

window.onload = main
