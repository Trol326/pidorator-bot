package roles

import "github.com/bwmarrin/discordgo"

func CreateRole(discord *discordgo.Session, guildID, name string) (*discordgo.Role, error) {
	params := &discordgo.RoleParams{
		Name: name,
	}

	role, err := discord.GuildRoleCreate(guildID, params)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func GetRole(discord *discordgo.Session, guildID, roleID string) (*discordgo.Role, error) {
	role, err := discord.State.Role(guildID, roleID)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func SetUserRole(discord *discordgo.Session, guildID, userID, roleID string) error {
	err := discord.GuildMemberRoleAdd(guildID, userID, roleID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserRole(discord *discordgo.Session, guildID, userID, roleID string) error {
	err := discord.GuildMemberRoleRemove(guildID, userID, roleID)
	if err != nil {
		return err
	}
	return nil
}

func MemberHasPermission(s *discordgo.Session, guildID string, userID string, permission int64) (bool, error) {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		guild, err = s.Guild(guildID)
		if err != nil {
			return false, err
		}
	}
	if guild.OwnerID == userID {
		return true, nil
	}

	member, err := s.State.Member(guildID, userID)
	if err != nil {
		if member, err = s.GuildMember(guildID, userID); err != nil {
			return false, err
		}
	}

	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			return false, err
		}
		if role.Permissions&permission != 0 {
			return true, nil
		}
	}

	return false, nil
}
