pub contract UserProfiles {
    pub struct HeaderPictures {
        pub(set) var smallUrl: String
        pub(set) var mediumUrl: String
        pub(set) var bigUrl: String

        init() {
            self.smallUrl = ""
            self.mediumUrl = ""
            self.bigUrl = ""
        }
    }

    priv let usernames: {Address:String}
    priv let avatars: {Address:{String:String}}
    priv let headerPics: {Address: HeaderPictures}

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

    pub fun getAllAvatarNames(_ addr: Address): [String] {
        return self.getAllAvatars(addr).keys
    }

    pub fun getAvatarByName(_ addr: Address, _ name: String): String? {
        return self.getAllAvatars(addr)[name]
    }

    pub fun getHeaderPics(_ addr: Address): HeaderPictures? {
        return self.headerPics[addr]
    }

    pub fun setHeaderPics(_ acc: AuthAccount, _ smallUrl: String, _ mediumUrl: String, _ bigUrl: String) {
        let headerPics = self.getHeaderPics(acc.address) ?? HeaderPictures()
        headerPics.bigUrl = bigUrl
        headerPics.mediumUrl = mediumUrl
        headerPics.smallUrl = smallUrl
        self.headerPics[acc.address] = headerPics;
    }

    init() {
        self.usernames = {}
        self.avatars = {}
        self.headerPics = {}
    }
}
