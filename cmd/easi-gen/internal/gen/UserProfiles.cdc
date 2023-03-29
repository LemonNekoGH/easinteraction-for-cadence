pub contract UserProfiles {
    access(self) let usernames: {Address:String}
    access(self) let avatars: {Address:{String:String}}

    pub fun setName(user acc: AuthAccount, to name: String) {
        self.usernames[acc.address] = name
    }

    pub fun getName(_ addr: Address): String {
        return self.usernames[addr] ?? ""
    }

    pub fun setAvatar(_ avatarName: String, _ avatarUrl: String, _ acc: AuthAccount) {
        let avatars = self.getAllAvatars(acc.address)
        avatars[avatarName] = avatarUrl
        self.avatars[acc.address] = avatars
    }

    pub fun getAllAvatars(_ addr: Address): {String:String} {
        return self.avatars[addr] ?? {}
    }

    init() {
        self.usernames = {}
        self.avatars = {}
    }
}
