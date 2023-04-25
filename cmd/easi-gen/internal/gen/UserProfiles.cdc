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

    pub resource SocialLink {
        pub var addr: Address
        pub var type: String
        pub var link: String

        init(_ addr: Address, _ type: String, _ link: String) {
            self.addr = addr
            self.type = type
            self.link = link
        }
    }

    pub event UsernameUpdate(_ name: String)
    pub event SocialLinkUpdate(_ user: Address, _ type: String, _ value: String)

    priv let usernames: {Address:String}
    priv let avatars: {Address:{String:String}}
    priv let headerPics: {Address: HeaderPictures}
    priv let socialLinks: @[SocialLink]

    pub fun setName(user acc: AuthAccount, to name: String) {
        self.usernames[acc.address] = name
        emit UsernameUpdate(name)
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

    pub fun setSocialLink(_ addr: Address, _ type: String, _ link: String) {
        var newOne <- create SocialLink(addr, type, link)
        self.setSocialLinkDirect(<- newOne)
    }

    pub fun getSocialLink(_ addr: Address, _ type: String): &SocialLink? {
        var index = 0
        // check if exists
        var found = false
        while (index < self.socialLinks.length) {
            if (self.socialLinks[index].type == type) {
                return &self.socialLinks[index] as &SocialLink?
            }
        }
        return nil
    }

    pub fun createSocialLink(_ addr: Address, _ type: String, _ link: String): @SocialLink {
        return <- create SocialLink(addr, type, link)
    }

    pub fun setSocialLinkDirect(_ link: @SocialLink) {
        var index = 0
        // check if exists
        while (index < self.socialLinks.length) {
            if (self.socialLinks[index].type != link.type) {
                continue
            }
            // exists, remove
            var found <- self.socialLinks.remove(at: index)
            destroy found
            index = index + 1
        }
        // add new
        emit SocialLinkUpdate(link.addr, link.type, link.link)
        self.socialLinks.append(<- link)
    }

    init() {
        self.usernames = {}
        self.avatars = {}
        self.headerPics = {}
        self.socialLinks <- []
    }
}
