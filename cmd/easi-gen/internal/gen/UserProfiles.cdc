pub contract UserProfiles {
    access(self) let usernames: {Address:String}

    pub fun setName(_ acc: AuthAccount, _ name: String) {
        self.usernames[acc.address] = name
    }

    pub fun getName(_ addr: Address): String {
        return self.usernames[addr] ?? ""
    }

    init() {
        self.usernames = {}
    }
}
