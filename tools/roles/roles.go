package roles

import "github.com/bwmarrin/discordgo"

func CreateRole(discord *discordgo.Session, guildID string) (*discordgo.Role, error) {
	params := &discordgo.RoleParams{
		Name: "Пидор дня",
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
